package client

import (
	"errors"
	"fmt"
	"strings"

	"gitar/pkg/client/common"
	"gitar/pkg/client/gitee"
	"gitar/pkg/client/github"
	"gitar/pkg/config"
)

func ParseRepoUrl(url string) (*common.RepoUrl, error) {
	if len(url) <= 0 {
		return nil, errors.New("url is empty")
	}
	if strings.Contains(url, github.Host) {
		return github.ParseGithubRepoUrl(url)
	}
	if strings.Contains(url, gitee.Host) {
		return gitee.ParseGiteeRepoUrl(url)
	}
	return nil, fmt.Errorf("unsupported url %s", url)
}

func ResolveArchive(url common.RepoUrl, config *config.ConfigProperties) (*common.ArchiveInfo, error) {
	if url.Platform == github.Platform {
		svc := github.NewGitHubService(config.GitHub.Token)
		return svc.ResolveArchive(url)
	}
	return nil, fmt.Errorf("unsupported platform: %s", url.Platform)
}
