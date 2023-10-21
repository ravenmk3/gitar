package utils

import (
	"os"
)

func FileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if err == nil {
		return !info.IsDir(), nil
	}
	return !os.IsNotExist(err), nil
}
