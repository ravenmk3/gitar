package git

import (
	"fmt"
	"regexp"
)

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
