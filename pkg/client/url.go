package client

import (
	"errors"
	"fmt"
	"strings"
)

const (
	Github     = "github"
	Gitee      = "gitee"
	GithubHost = "github.com"
	GiteeHost  = "gitee.com"
)

type RepoUrl struct {
	Platform string
	Host     string
	Owner    string
	Repo     string
	Release  string
	Tag      string
	Branch   string
	Commit   string
	RefName  string
}

type ArchiveUrl struct {
	Platform string
	Commit   string
	Url      string
}

func ParseRepoUrl(url string) (*RepoUrl, error) {
	if len(url) <= 0 {
		return nil, errors.New("url is empty")
	}
	if strings.Contains(url, GithubHost) {
		return ParseGithubRepoUrl(url)
	}
	if strings.Contains(url, GiteeHost) {
		return ParseGiteeRepoUrl(url)
	}
	return nil, fmt.Errorf("unsupported url %s", url)
}
