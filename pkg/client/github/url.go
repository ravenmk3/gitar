package github

import (
	"fmt"
	"net/url"
	"regexp"

	"gitar/pkg/client/common"
)

func ParseGithubRepoUrl(rawUrl string) (*common.RepoUrl, error) {
	info := &common.RepoUrl{
		Platform: Platform,
		Host:     Host,
	}

	// SSH URL
	re := regexp.MustCompile(`^git@github\.com:([\w\-.]+)/([\w\-.]+)\.git$`)
	match := re.FindStringSubmatch(rawUrl)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		return info, nil
	}

	// HTTPS URL
	re = regexp.MustCompile(`^https://github\.com/([\w\-.]+)/([\w\-.]+)(?:\.git)?$`)
	match = re.FindStringSubmatch(rawUrl)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		return info, nil
	}

	// HTTPS Tag URL
	re = regexp.MustCompile(`^https://github\.com/([\w\-.]+)/([\w\-.]+)/releases/tag/([\w\-.%]+)?$`)
	match = re.FindStringSubmatch(rawUrl)
	if match != nil {
		tag, err := url.QueryUnescape(match[3])
		if err != nil {
			return nil, err
		}

		info.Owner = match[1]
		info.Repo = match[2]
		info.Release = tag
		info.Tag = tag
		return info, nil
	}

	// HTTPS Tree of Commit URL
	re = regexp.MustCompile(`^https://github\.com/([\w\-.]+)/([\w\-.]+)/tree/([0-9a-fA-F]{40})?$`)
	match = re.FindStringSubmatch(rawUrl)
	if match != nil {
		info.Owner = match[1]
		info.Repo = match[2]
		info.Commit = match[3]
		return info, nil
	}

	// HTTPS Tree URL
	re = regexp.MustCompile(`^https://github\.com/([\w\-.]+)/([\w\-.]+)/tree/([\w\-.%]+)?$`)
	match = re.FindStringSubmatch(rawUrl)
	if match != nil {
		ref, err := url.QueryUnescape(match[3])
		if err != nil {
			return nil, err
		}
		info.Owner = match[1]
		info.Repo = match[2]
		info.Branch = ref
		info.RefName = ref
		return info, nil
	}

	return nil, fmt.Errorf("unsupported GitHub url %s", rawUrl)
}
