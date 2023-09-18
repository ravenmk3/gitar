package git

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

type GitUrl struct {
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

func ParseGitUrl(url string) (*GitUrl, error) {
	if len(url) <= 0 {
		return nil, errors.New("url is empty")
	}
	if strings.Contains(url, GithubHost) {
		return ParseGithubUrl(url)
	}
	if strings.Contains(url, GiteeHost) {
		return ParseGiteeUrl(url)
	}
	return nil, fmt.Errorf("unsupported url %s", url)
}
