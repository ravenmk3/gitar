package git

import (
	"errors"
	"fmt"
	"regexp"
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

func ParseGithubUrl(url string) (*GitUrl, error) {
	info := &GitUrl{
		Platform: Github,
		Host:     GithubHost,
	}

	// SSH URL
	re := regexp.MustCompile(`^git@github\.com:([\w\-.]+)/([\w\-.]+)\.git$`)
	match := re.FindStringSubmatch(url)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		return info, nil
	}

	// HTTPS URL
	re = regexp.MustCompile(`^https://github\.com/([\w\-.]+)/([\w\-.]+)(?:\.git)?$`)
	match = re.FindStringSubmatch(url)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		return info, nil
	}

	// HTTPS Tag URL
	re = regexp.MustCompile(`^https://github\.com/([\w\-.]+)/([\w\-.]+)/releases/tag/([\w\-.]+)?$`)
	match = re.FindStringSubmatch(url)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		info.Release = match[3]
		info.Tag = match[3]
		return info, nil
	}

	// HTTPS Tree of Commit URL
	re = regexp.MustCompile(`^https://github\.com/([\w\-.]+)/([\w\-.]+)/tree/([0-9a-fA-F]{40})?$`)
	match = re.FindStringSubmatch(url)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		info.Commit = match[3]
		return info, nil
	}

	// HTTPS Tree URL
	re = regexp.MustCompile(`^https://github\.com/([\w\-.]+)/([\w\-.]+)/tree/([\w\-.]+)?$`)
	match = re.FindStringSubmatch(url)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		info.RefName = match[3]
		return info, nil
	}

	return nil, fmt.Errorf("unsupported GitHub url %s", url)
}

func ParseGiteeUrl(url string) (*GitUrl, error) {
	info := &GitUrl{
		Platform: Gitee,
		Host:     GiteeHost,
	}

	// SSH URL
	re := regexp.MustCompile(`^git@gitee\.com:([\w\-.]+)/([\w\-.]+)\.git$`)
	match := re.FindStringSubmatch(url)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		return info, nil
	}

	// HTTPS URL
	re = regexp.MustCompile(`^https://gitee\.com/([\w\-.]+)/([\w\-.]+)(?:\.git)?$`)
	match = re.FindStringSubmatch(url)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		return info, nil
	}

	// HTTPS Tag URL
	re = regexp.MustCompile(`^https://gitee\.com/([\w\-.]+)/([\w\-.]+)/releases/tag/([\w\-.]+)?$`)
	match = re.FindStringSubmatch(url)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		info.Release = match[3]
		info.Tag = match[3]
		return info, nil
	}

	// HTTPS Tree of Commit URL
	re = regexp.MustCompile(`^https://gitee\.com/([\w\-.]+)/([\w\-.]+)/tree/([0-9a-fA-F]{40})?$`)
	match = re.FindStringSubmatch(url)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		info.Commit = match[3]
		return info, nil
	}

	// HTTPS Tree URL
	re = regexp.MustCompile(`^https://gitee\.com/([\w\-.]+)/([\w\-.]+)/tree/([\w\-.]+)?$`)
	match = re.FindStringSubmatch(url)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		info.RefName = match[3]
		return info, nil
	}

	return nil, fmt.Errorf("unsupported Gitee url %s", url)
}
