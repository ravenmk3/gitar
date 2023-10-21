package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gitar/pkg/client"
	"gitar/pkg/config"
	"gitar/pkg/data"
	"gitar/pkg/utils"
	"github.com/sirupsen/logrus"
)

func DownloadArchive(url string, shouldSendMail bool) error {
	err := DoDownloadArchive(url, shouldSendMail)
	logrus.Infof("All done")
	return err
}

func DoDownloadArchive(url string, shouldSendMail bool) error {
	logrus.Infof("Download archive")

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	logrus.Infof("Config: %+v", *cfg)

	logrus.Infof("Git URL: %s", url)
	repoUrl, err := client.ParseRepoUrl(url)
	if err != nil {
		return err
	}
	logrus.Infof("Parsed URL: %+v", *repoUrl)

	arc, err := client.ResolveArchive(*repoUrl)
	if err != nil {
		return err
	}
	logrus.Infof("Archive: %+v", *arc)

	err = os.MkdirAll(cfg.Paths.Data, os.ModePerm)
	if err != nil {
		return err
	}

	store := data.NewSqlite3DataStore(filepath.Join(cfg.Paths.Data, "gitar.sqlite"))
	err = store.Open()
	if err != nil {
		return err
	}

	repoKey := fmt.Sprintf("%s:%s/%s", repoUrl.Platform, repoUrl.Owner, repoUrl.Repo)
	err = store.SaveRepo(repoKey)
	if err != nil {
		return err
	}

	markDownloaded, err := store.IsCommitDownloaded(arc.Commit)
	if err != nil {
		return err
	}

	if markDownloaded && !shouldSendMail {
		logrus.Warnf("Commit already downloaded: %s", arc.Commit)
		return nil
	}

	arcFile := fmt.Sprintf("%s.tar.xz", arc.Name)
	destDir := filepath.Join(cfg.Paths.Repo, repoUrl.Platform, repoUrl.Owner, repoUrl.Repo)
	destPath := filepath.Join(destDir, arcFile)

	tempFile := fmt.Sprintf("%s-%s.tar.gz", arc.Name, arc.Commit)
	tempPath := filepath.Join(cfg.Paths.Temp, tempFile)
	tempXzFile := fmt.Sprintf("%s-%s.tar.xz", arc.Name, arc.Commit)
	tempXzPath := filepath.Join(cfg.Paths.Temp, tempXzFile)
	logrus.Infof("Temp file: %s", tempPath)

	destExists, err := utils.FileExists(destPath)
	if err != nil {
		return err
	}

	if destExists {
		logrus.Warnf("Already downloaded: %s", destPath)
	} else {
		err = os.RemoveAll(tempPath)
		if err != nil {
			return err
		}

		err = utils.Aria2Download(arc.TarUrl, cfg.Paths.Temp, tempFile)
		if err != nil {
			return err
		}

		logrus.Infof("Converting gzip archive to xz")
		err = utils.Gzip2Xz(tempPath, tempXzPath)
		if err != nil {
			return err
		}

		err = os.RemoveAll(tempPath)
		if err != nil {
			return err
		}

		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			return err
		}

		logrus.Infof("Downloaded: %s", destPath)
		err = os.Rename(tempXzPath, destPath)
		if err != nil {
			return err
		}
	}

	if !markDownloaded {
		err = store.SetCommitDownloaded(arc.Commit)
	}

	if shouldSendMail {

		mailed, err := store.IsCommitMailed(arc.Commit)
		if err != nil {
			return err
		}

		if mailed {
			logrus.Warnf("Commit already mailed: %s", arc.Commit)
		} else {
			err := sendMailWithRetry(destPath, 999)
			if err != nil {
				return err
			}

			err = store.SetCommitMailed(arc.Commit)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func sendMailWithRetry(file string, maxAttempts int) error {
	for i := 0; i < maxAttempts; i++ {
		err := sendMail(file)
		if err != nil {
			logrus.Error(err)
			time.Sleep(calcMailRetryDelay(i + 1))
			continue
		}
		return nil
	}
	return errors.New("send mail failed")
}

func sendMail(file string) error {
	cmd := exec.Command("filemailer", "send", "--profile=gitar", file)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

func calcMailRetryDelay(num int) time.Duration {
	delay := num * num * num
	if delay > 7200 {
		delay = 7200
	}
	return time.Second * time.Duration(delay)
}
