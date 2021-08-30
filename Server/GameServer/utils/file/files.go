package file

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetExecutePath() (string, error) {
	dir, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(dir)
	return exPath, nil
}

func GetWorkPath() (string, error) {

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

func PathConversionByOs(path string) string {
	switch runtime.GOOS {
	case "darwin","linux":
		return strings.ReplaceAll(path,"\\","/")
	default: //"windows"
		return strings.ReplaceAll(path,"/","\\")
	}
}
