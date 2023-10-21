package app

import (
	"fmt"
	"os"
	"path/filepath"

	"gitar/pkg/client"
	"gitar/pkg/config"
	"gitar/pkg/utils"
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
	logrus.Infof("Parsed URL: %+v", *repoUrl)

	arc, err := client.ResolveArchive(*repoUrl)
	if err != nil {
		return err
	}
	logrus.Infof("Archive: %+v", *arc)

	arcFile := fmt.Sprintf("%s.tar.gz", arc.Name)
	destDir := filepath.Join(cfg.RepoDir, repoUrl.Platform, repoUrl.Owner, repoUrl.Repo)
	destPath := filepath.Join(destDir, arcFile)
	// TODO 判断是否已存在

	tempFile := fmt.Sprintf("%s-%s.tar.gz", arc.Name, arc.Commit)
	tempPath := filepath.Join(cfg.TempDir, tempFile)
	// TODO 删除旧文件
	err = utils.Aria2Download(arc.TarUrl, cfg.TempDir, tempFile)
	if err != nil {
		return err
	}
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	logrus.Infof("Move file %s => %s", tempPath, destPath)
	err = os.Rename(tempPath, destPath)
	if err != nil {
		return err
	}

	return err
}
