package util

import (
	"io"
	"io/ioutil"
	"os"
)

var FileHelper = &fileHelper{}

type fileHelper struct {
}

func (f *fileHelper) CreateDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, 0755)
	}
}

func (f *fileHelper) ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}

// 复制文件
func (f *fileHelper) CopyFile(sourcePath string, destPath string) error {
	srcFile, err := os.Open(sourcePath)

	if err != nil {
		return err
	}

	defer srcFile.Close()

	// 打开dstFileName
	dstFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)

	return err

}
