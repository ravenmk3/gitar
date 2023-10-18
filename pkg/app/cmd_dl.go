package app

import (
	"gitar/pkg/client"
	"gitar/pkg/client/github"
	"gitar/pkg/config"
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
	repoUrl, err := client.ParseRepoUrl(url)
	if err != nil {
		return err
	}

	logrus.Infof("%+v", repoUrl)

	arc, err := github.ResolveGithubArchive(repoUrl)
	logrus.Infof("%+v", arc)
	return err
}
