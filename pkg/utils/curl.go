package utils

import (
	"os"
	"os/exec"
)

func CurlDownload(url string, dir string, file string, maxTries int) error {
	if maxTries < 0 {
		maxTries = 99999999
	} else if maxTries < 1 {
		maxTries = 1
	}

	args := []string{
		"--location",
		"--output",
		file,
		url,
	}

	var err error
	for i := 0; i < maxTries; i++ {
		err = execCurl(dir, args, true)
		if err == nil {
			return nil
		}
	}

	return err
}

func execCurl(dir string, args []string, redirect bool) error {
	cmd := exec.Command("curl", args...)
	cmd.Dir = dir
	if redirect {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
