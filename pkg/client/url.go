package client

import (
	"errors"
	"fmt"
	"strings"

	"gitar/pkg/client/common"
	"gitar/pkg/client/gitee"
	"gitar/pkg/client/github"
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
