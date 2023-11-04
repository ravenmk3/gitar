package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func Aria2Download(url string, dir string, file string, maxTries int) error {
	if maxTries < 0 {
		maxTries = 99999
	} else if maxTries < 1 {
		maxTries = 1
	}
	args := []string{
		"--dir=" + dir,
		"--out=" + file,
		fmt.Sprintf("--max-tries=%d", maxTries),
		"--lowest-speed-limit=1",
		"--user-agent=" + HttpUserAgent,
		"--file-allocation=none",
		url,
	}
	return execAria2c(args, true)
}

func execAria2c(args []string, redirect bool) error {
	cmd := exec.Command("aria2c", args...)
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
