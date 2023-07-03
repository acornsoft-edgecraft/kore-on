package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"kore-on/pkg/config"
	"kore-on/pkg/logger"
	"kore-on/pkg/utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"kore-on/cmd/koreonctl/conf"

	"github.com/elastic/go-sysinfo"
	"github.com/spf13/cobra"
)

type strInitCmd struct {
	verbose        bool
	osRelease      string
	osArchitecture string
}

func initCmd() *cobra.Command {
	init := &strInitCmd{}
	cmd := &cobra.Command{
		Use:          "init [flags]",
		Short:        "Get configuration file",
		Long:         "This command downloads a sample file that can set machine information and installation information.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return init.run()
		},
	}

	// SubCommand add
	cmd.AddCommand(emptyCmd())

	// SubCommand validation
	utils.CheckCommand(cmd)

	f := cmd.Flags()
	f.BoolVar(&init.verbose, "vvv", false, "verbose")

	return cmd
}

func (c *strInitCmd) run() error {
	// work directory tree check
	workDir, err := checkDirTree()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// Check installed Podman
	if err := installPodman(workDir); err != nil {
		logger.Fatal(err)
	}

	// system info
	host, err := sysinfo.Host()
	if err != nil {
		logger.Fatal(err)
	}
	c.osArchitecture = host.Info().Architecture
	c.osRelease = host.Info().OS.Platform

	logger.Infof("Start provisioning for cloud infrastructure")

	if err := c.init(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strInitCmd) init(workDir string) error {

	currTime := time.Now()

	SUCCESS_FORMAT := "\033[1;32m%s\033[0m\n"
	koreOnConfigFile := conf.KoreOnConfigFile

	infoStr := "Do you really want to init?\n" +
		"Is this ok [y/n]: "
	if !utils.CheckUserInput(infoStr, "y") {
		fmt.Println("nothing to changed. exit")
		os.Exit(1)
	}

	koreOnConfigFilePath, err := filepath.Abs(koreOnConfigFile)
	if err != nil {
		ioutil.WriteFile(workDir+"/"+koreOnConfigFile, []byte(config.Template), 0600)
		fmt.Printf(SUCCESS_FORMAT, fmt.Sprintf("Initialize completed, Edit %s file according to your environment and run `koreonctl create`", koreOnConfigFile))
	} else {
		fmt.Println("Previous " + koreOnConfigFile + " file exist and it will be backup")
		os.Rename(koreOnConfigFilePath, koreOnConfigFilePath+"_"+currTime.Format("20060102150405"))
		ioutil.WriteFile(workDir+"/"+koreOnConfigFile, []byte(config.Template), 0600)
		fmt.Printf(SUCCESS_FORMAT, fmt.Sprintf("Initialize completed, Edit %s file according to your environment and run `koreonctl create`", koreOnConfigFile))
	}

	return nil
}

func installPodman(workDir string) error {
	// podmand installed check
	_, podmanCheck := exec.LookPath("podman")
	if podmanCheck == nil {
		logger.Info("podman already.")
		return nil
	}

	if runtime.GOOS != "linux" {
		errStr := "Installation of the podman package is only supported on Linux platforms.\n" +
			"If your system is not running on a Linux platform, please manually install the podman package and then run it again."

		return fmt.Errorf("%s", errStr)
	}

	infoStr := "Installing Podman is mandatory. Do you want to proceed with the Podman installation? If not, please install it manually. Once the installation is complete, please run it again.\n" +
		"Is this ok [y/n]: "
	if !utils.CheckUserInput(infoStr, "y") {
		return fmt.Errorf("nothing to changed. exit")
	}

	// tar.gz 압축 파일 열기
	file, err := os.Open(workDir + "/archive/podman-linux-amd64.tar.gz")
	if err != nil {
		return err
	}
	defer file.Close()

	// gzip 해제
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// tar 아카이브 열기
	tarReader := tar.NewReader(gzipReader)

	excludePath := "/README.md"

	// 압축 해제된 파일들을 시스템에 푸는 작업 수행
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// 아카이브의 끝에 도달하면 반복문 종료
			break
		}
		if err != nil {
			return err
		}

		// 풀어질 파일의 경로 생성 / 특정 경로 제거
		subPath := strings.Split(header.Name, "/")
		targetPath := removePath(header.Name, subPath[0])

		// 특정 파일 또는 디렉토리를 제외합니다.
		if targetPath == excludePath || isDescendant(excludePath, targetPath) {
			continue
		}

		// 파일 또는 디렉토리 생성
		if header.Typeflag == tar.TypeDir {
			// 디렉토리 생성
			err := os.MkdirAll(targetPath, 0755)
			if err != nil {
				return err
			}
		} else if header.Typeflag == tar.TypeReg {
			// 파일 생성
			file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer file.Close()

			// 파일 내용을 복사
			_, err = io.Copy(file, tarReader)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func removePath(path, subPath string) string {
	// 특정 경로 제거
	result := strings.Replace(path, subPath, "", 1)

	// 최종 경로 정리
	result = filepath.Clean(result)

	return result
}

// 경로가 제외 경로의 하위 경로인지 확인하는 함수
func isDescendant(excludePath, path string) bool {
	relative, err := filepath.Rel(excludePath, path)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(relative, "..")
}

func checkDirTree() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// 실행 파일이 있는 디렉토리 경로 추출
	executableDir := filepath.Dir(executablePath)

	// 실행 파일의 상위 경로 추출
	parentDir := filepath.Dir(executableDir)

	dirTree := []string{
		parentDir + "/bin",
		parentDir + "/archive",
		parentDir + "/config",
		parentDir + "/logs",
	}
	currentDirTree := []string{}
	err = filepath.WalkDir(parentDir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Directory
			// Check if the string contains a specific substring
			for _, v := range dirTree {
				if v == path {
					currentDirTree = append(currentDirTree, v)
				}
			}
		} else {
			// File
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if len(dirTree) != len(currentDirTree) {
		errStr := ""
		for _, dirValue := range dirTree {
			found := false
			for _, currentValue := range currentDirTree {
				if dirValue == currentValue {
					found = true
					break
				}
			}
			if !found {
				errStr = errStr + fmt.Sprintf("%s No such directory ", dirValue)
			}
		}
		return "", fmt.Errorf(errStr)
	}
	return parentDir, nil
}
