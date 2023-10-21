package utils

import (
	"os"
	"os/exec"
)

func Aria2Download(url string, dir string, file string) error {
	args := []string{
		"--dir=" + dir,
		"--out=" + file,
		"--max-tries=10",
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
