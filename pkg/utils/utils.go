package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kore-on/pkg/logger"
	"os"
	"sort"
	"strings"

	"github.com/spf13/viper"
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

	if len(strings.Split(version, ".")) == 2 {
		k := version
		v, err := supportversion[version]
		if !err {
			return ""
		}
		version = fmt.Sprintf("%v.%v", k, v[len(v)-1])
	}

	keys := make([]string, 0, len(supportversion))
	for k := range supportversion {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	values, err := supportversion[keys[0]]
	if !err {
		return ""
	}

	latest := fmt.Sprintf("%v.%v", keys[0], values[len(values)-1])

	if version == "" {
		return latest
	} else {
		major := version[0:strings.LastIndex(version, ".")]
		minor := version[len(major)+1 : len(version)+0]
		for _, v := range supportversion[major] {
			if v == minor {
				return version
			}
		}
		// Returns no matching version
		return ""
	}
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

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
