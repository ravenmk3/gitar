package utils

import (
	"io"
	"os"
	"runtime"
)

func FileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if err == nil {
		return !info.IsDir(), nil
	}
	return !os.IsNotExist(err), nil
}

func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func(srcFile *os.File) {
		_ = srcFile.Close()
	}(srcFile)

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func(dstFile *os.File) {
		_ = dstFile.Close()
	}(dstFile)

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func MoveFile(srcPath, dstPath string) error {
	if runtime.GOOS == "windows" {
		return os.Rename(srcPath, dstPath)
	}

	err := CopyFile(srcPath, dstPath)
	if err != nil {
		return err
	}
	return os.RemoveAll(srcPath)
}
