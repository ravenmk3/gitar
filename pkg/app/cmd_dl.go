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
	"gitar/pkg/fslock"
	"gitar/pkg/utils"
	"github.com/sirupsen/logrus"
)

func DownloadArchive(url string, shouldSendMail bool) error {
	err := DoDownloadArchive(url, shouldSendMail)
	if err == nil {
		logrus.Infof("All done")
	}
	return err
}

func DoDownloadArchive(url string, shouldSendMail bool) error {
	logrus.Infof("Downloading archive")

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	logrus.Infof("Config: %+v", *cfg)

	logrus.Infof("URL: %s", url)
	repoUrl, err := client.ParseRepoUrl(url)
	if err != nil {
		return err
	}
	logrus.Infof("Platform: %s", repoUrl.Platform)
	logrus.Infof("Repository: %s/%s", repoUrl.Owner, repoUrl.Repo)
	logrus.Infof("Parsed-Tag: %s", repoUrl.Tag)
	logrus.Infof("Parsed-Branch: %s", repoUrl.Branch)
	logrus.Infof("Parsed-Commit: %s", repoUrl.Commit)

	arc, err := client.ResolveArchive(*repoUrl, cfg.Token)
	if err != nil {
		return err
	}
	logrus.Infof("Archive-Name: %s", arc.Name)
	logrus.Infof("Archive-Commit: %s", arc.Commit)
	logrus.Infof("Archive-Tar: %s", arc.TarUrl)
	logrus.Infof("Archive-Zip: %s", arc.ZipUrl)

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
		lockFile := filepath.Join(cfg.Paths.Temp, arc.Commit+".lock")
		lock := fslock.New(lockFile)
		err := lock.TryLock()
		if err != nil {
			return err
		}
		defer func(lock fslock.Lock) {
			err := lock.Unlock()
			if err != nil {
				logrus.Error(err)
			}
		}(lock)

		err = os.RemoveAll(tempPath)
		if err != nil {
			return err
		}

		err = utils.Aria2Download(arc.TarUrl, cfg.Paths.Temp, tempFile)
		if err != nil {
			return err
		}

		gzipSize, err := utils.GetFileSize(tempPath)
		if err != nil {
			return err
		}

		logrus.Infof("Downloaded: %s (%s)", tempFile, utils.HumanReadableSize(gzipSize))
		logrus.Infof("Converting gzip archive to xz")
		err = utils.Gzip2Xz(tempPath, tempXzPath)
		if err != nil {
			return err
		}

		xzSize, err := utils.GetFileSize(tempXzPath)
		if err != nil {
			return err
		}
		logrus.Infof("Converted: gzip (%s) => xz (%s)",
			utils.HumanReadableSize(gzipSize),
			utils.HumanReadableSize(xzSize))

		err = os.RemoveAll(tempPath)
		if err != nil {
			return err
		}

		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			return err
		}

		err = utils.MoveFile(tempXzPath, destPath)
		if err != nil {
			return err
		}
		logrus.Infof("Saved: %s", destPath)
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
			subject := fmt.Sprintf("%s:%s/%s.tar.xz", repoUrl.Platform, repoUrl.Owner, arc.Name)

			err := sendMailWithRetry(destPath, subject, 999)
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

func sendMailWithRetry(file, subject string, maxAttempts int) error {
	logrus.Infof("Sending email")
	for i := 0; i < maxAttempts; i++ {
		err := sendMail(file, subject)
		if err != nil {
			logrus.Error(err)
			delay := calcMailRetryDelay(i + 1)
			logrus.Infof("Retry after %s", delay)
			time.Sleep(delay)
			continue
		}
		return nil
	}
	return errors.New("send mail failed")
}

func sendMail(file string, subject string) error {
	cmd := exec.Command("filemailer", "send", "--profile=gitar", "--subject", subject, file)
	cmd.Stdout = os.Stdout
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
