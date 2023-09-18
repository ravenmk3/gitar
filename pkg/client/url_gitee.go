package client

import (
	"fmt"
	"regexp"
)

func ParseGiteeRepoUrl(url string) (*RepoUrl, error) {
	info := &RepoUrl{
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