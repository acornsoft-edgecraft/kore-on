package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"kore-on/pkg/logger"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"gopkg.in/yaml.v2"
)

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func FileSizeAndExists(name string) (int64, bool, error) {
	//workDir, _ := os.Getwd()
	//fmt.Printf("workdir %s\n",workDir)
	var size int64 = 0
	stat, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return size, false, err
		}
	}
	size = stat.Size()

	return size, true, nil
}

func CopyFile(source string, dest string) (err error) {
	//fmt.Printf("copy file source %s dest %s \n",source, dest)
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		fmt.Printf("error %s\n", err)
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)

		if err == nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		} else {
			return err
		}
	} else {
		return err
	}

	return nil
}

func CopyFile0600(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		err = os.Chmod(dest, 0600)
	} else {
		return err
	}

	return nil
}

func IskoreOnConfigFilePath(s string) string {
	currDir, _ := os.Getwd()
	sub := viper.GetString("KoreOn.KoreOnConfigFileSubDir")
	if sub != "" {
		sub = "/" + sub + "/"
	} else {
		sub = "/"
	}
	return currDir + sub + s
}

func IsSupportVersion(version string, conf string) string {
	supportversion := viper.GetStringMapStringSlice(conf)
	if len(supportversion) == 0 {
		logger.Fatal("koreon > There is no supported version.")
	}
	if len(strings.Split(version, ".")) == 2 {
		k := version
		v, err := supportversion[version]
		if !err {
			logger.Fatal("koreon > There is no supported version.")
		}
		if len(v) == 1 && v[0] == "" {
			version = fmt.Sprintf("%v", k)
		} else {
			version = fmt.Sprintf("%v.%v", k, v[len(v)-1])
		}
	}

	keys := make([]string, 0, len(supportversion))
	for k := range supportversion {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	values, err := supportversion[keys[0]]
	if !err {
		logger.Fatal("koreon > There is no supported version.")
	}

	latest := fmt.Sprintf("%v.%v", keys[0], values[len(values)-1])

	if version == "" || version == "latest" {
		return latest
	} else {
		major := version[0:strings.LastIndex(version, ".")]
		minor := version[len(major)+1 : len(version)+0]

		for _, v := range supportversion[major] {
			if v == minor {
				return version
			}
		}
		// Returns just use major version
		return version
	}
}

func GetSupportVersion(version string, key string) map[string]interface{} {
	getVersion := viper.GetStringMap("SupportVersion")
	chekVersion := false

	for k, v := range getVersion[key].(map[string]interface{}) {
		if k == version {
			chekVersion = true
			return v.(map[string]interface{})
		}
	}
	if !chekVersion {
		for k, v := range getVersion[key].(map[string]interface{}) {
			if k == version[0:strings.LastIndex(version, ".")] {
				return v.(map[string]interface{})
			}
		}
	}

	return nil
}

func ListSupportVersion(conf string) string {
	supportversion := viper.GetStringMapStringSlice(conf)

	keys := make([]string, 0, len(supportversion))
	for k := range supportversion {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	for k, v := range supportversion {
		for i, j := range v {
			supportversion[k][i] = fmt.Sprintf("%v.%v", k, j)
		}
	}

	b, err := json.Marshal(supportversion)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	b, _ = prettyPrint(b)
	return string(b)
}

func CheckDocker() error {
	//fmt.Println("Checking pre-requisition [" + runtime.GOOS + "]")
	_, err := exec.Command("docker", "-v").Output()

	if err != nil {
		//fmt.Println(err.Error())
		logger.Fatal("docker is not found. Install docker before proceeding.")
		logger.Fatal("If it is a closed network, you can install it using the 'koreanctl bastion' command with the prepared package.")
		logger.Fatal("Visit https://www.docker.com/get-started")
		return err
	}
	return nil
}

func CheckUserInput(prompt string, checkWord string) bool {
	var res string
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	buf, _ := reader.ReadString('\n')

	if runtime.GOOS == "windows" {
		res = strings.Split(buf, "\r\n")[0]
	} else {
		res = strings.Split(buf, "\n")[0]
	}

	if res == checkWord {
		return true
	}

	return false
}

func ExecCommand(c string, commandArgs []string) *exec.Cmd {
	binary, lookErr := exec.LookPath(c)
	if lookErr != nil {
		logger.Fatal(lookErr)
	}

	cmd := exec.Command(binary, commandArgs...)
	return cmd
}

// InputPrompt receives a string value using the label
func InputPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

// The entered password will not be displayed on the screen
func SensitivePrompt(label string) string {
	var s string
	for {
		fmt.Fprint(os.Stderr, label+" ")
		pw, _ := term.ReadPassword(int(syscall.Stdin))
		s = string(pw)
		if s != "" {
			break
		} else {
			fmt.Println()
		}
	}
	fmt.Println()
	return s
}

func SetValuesFile(key string, v map[string]interface{}) (map[string]interface{}, error) {
	addonPath := viper.GetString("KoreOn.KoreOnConfigFileSubDir")
	var addonYaml string
	var dataYaml map[string]interface{}
	if v["ValuesFile"].(string) != "" {
		filename := filepath.Base(v["ValuesFile"].(string))
		addonYaml = addonPath + "/" + filename
		yamlFile, err := ioutil.ReadFile(addonYaml)
		if err != nil {
			logger.Fatal(err)
		}

		var values map[string]interface{}

		err = yaml.Unmarshal(yamlFile, &values)
		if err != nil {
			logger.Fatal(err)
		}

		valuesKey := strcase.ToSnake(key) + "_values"

		dataYaml = map[string]interface{}{
			valuesKey: values,
		}
	}

	return dataYaml, nil
}

// customTrasnformer
func OutputColored() results.TransformerFunc {
	return func(message string) string {
		yellow := color.New(color.FgYellow).SprintFunc()
		return fmt.Sprintf("%v", yellow(message))
	}
}
func StrucToJson(s interface{}) (map[string]interface{}, error) {
	var data map[string]interface{}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
func Print(b []byte) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return string(out.Bytes()), err
}
