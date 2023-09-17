package app

import (
	"gitar/pkg/config"
	"gitar/pkg/git"
	"github.com/sirupsen/logrus"
)

func DownloadArchive(url string) error {
	logrus.Infof("Download archive")

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	logrus.Infof("Config: %+v", *cfg)

	logrus.Infof("Git URL: %s", url)
	gitUrl, err := git.ParseGitUrl(url)
	if err != nil {
		return err
	}

	logrus.Infof("%+v", gitUrl)
	return err
}
