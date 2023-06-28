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

	workDir, _ := os.Getwd()
	var err error = nil
	logger.Infof("Start provisioning for cloud infrastructure")

	if err = c.init(workDir); err != nil {
		return err
	}
	return nil
}

func (c *strInitCmd) init(workDir string) error {
	currTime := time.Now()

	SUCCESS_FORMAT := "\033[1;32m%s\033[0m\n"
	koreOnConfigFile := "config/" + conf.KoreOnConfigFile

	infoStr := "Do you really want to init?\n" +
		"If you proceed, it will install the podman package on your system. The podman package is a mandatory requirement.\n" +
		"If you do not wish to install podman now, please manually install it and then run the process again!\n" +
		"Is this ok [y/n]: "
	if !utils.CheckUserInput(infoStr, "y") {
		fmt.Println("nothing to changed. exit")
		os.Exit(1)
	}

	koreOnConfigFilePath, _ := filepath.Abs(koreOnConfigFile)
	_, err := os.Stat(koreOnConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			ioutil.WriteFile(workDir+"/"+koreOnConfigFile, []byte(config.Template), 0600)
			fmt.Printf(SUCCESS_FORMAT, fmt.Sprintf("Initialize completed, Edit %s file according to your environment and run `koreonctl create`", koreOnConfigFile))
		} else {
			fmt.Println("Previous " + koreOnConfigFile + " file exist and it will be backup")
			os.Rename(koreOnConfigFilePath, koreOnConfigFilePath+"_"+currTime.Format("20060102150405"))
			ioutil.WriteFile(workDir+"/"+koreOnConfigFile, []byte(config.Template), 0600)
			fmt.Printf(SUCCESS_FORMAT, fmt.Sprintf("Initialize completed, Edit %s file according to your environment and run `koreonctl create`", koreOnConfigFile))
		}
	}

	c.installPodman(workDir)
	return nil
}

func (c *strInitCmd) installPodman(workDir string) error {
	// system info
	host, err := sysinfo.Host()
	if err != nil {
		logger.Fatal(err)
	}
	c.osArchitecture = host.Info().Architecture
	c.osRelease = host.Info().OS.Platform

	// podmand installed check
	_, podmanCheck := exec.LookPath("podman")
	if podmanCheck == nil {
		logger.Info("podman already.")
		os.Exit(1)
	}

	if runtime.GOOS != "linux" {
		errStr := "Installation of the podman package is only supported on Linux platforms." +
			"If your system is not running on a Linux platform, please manually install the podman package and then run it again."
		logger.Error(errStr)
		os.Exit(1)
	}

	// tar.gz 압축 파일 열기
	file, err := os.Open("./build/package/podman-linux-amd64.tar.gz")
	if err != nil {
		logger.Error("Error opening tar.gz file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// gzip 해제
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		logger.Error("Error creating gzip reader:", err)
		os.Exit(1)
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
			logger.Error("Error reading tar header:", err)
			os.Exit(1)
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
				logger.Error("Error creating directory:", err)
				os.Exit(1)
			}
		} else if header.Typeflag == tar.TypeReg {
			// 파일 생성
			file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				logger.Error("Error creating file:", err)
				os.Exit(1)
			}
			defer file.Close()

			// 파일 내용을 복사
			_, err = io.Copy(file, tarReader)
			if err != nil {
				logger.Error("Error extracting file contents:", err)
				os.Exit(1)
			}
		}

		logger.Info("Extracted file:", targetPath)
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
