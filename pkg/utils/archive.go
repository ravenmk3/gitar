package utils

import (
	"compress/gzip"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func Gzip2Xz(gzipPath, xzPath string) error {
	gzipFile, err := os.Open(gzipPath)
	if err != nil {
		return err
	}

	defer func(gzipFile *os.File) {
		err := gzipFile.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(gzipFile)

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return err
	}

	defer func(gzipReader *gzip.Reader) {
		err := gzipReader.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(gzipReader)

	xzFile, err := os.Create(xzPath)
	if err != nil {
		return err
	}

	defer func(xzFile *os.File) {
		err := xzFile.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(xzFile)

	xzCmd := exec.Command("xz", "-c", "-")
	xzCmd.Stdin = gzipReader
	xzCmd.Stdout = xzFile

	err = xzCmd.Start()
	if err != nil {
		return err
	}
	return xzCmd.Wait()
}
