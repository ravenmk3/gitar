package github

import (
	"errors"
	"fmt"
	"regexp"

	"gitar/pkg/client/common"
)

func ParseGithubRepoUrl(url string) (*common.RepoUrl, error) {
	info := &common.RepoUrl{
		Platform: Platform,
		Host:     Host,
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

func ResolveGithubArchive(url common.RepoUrl) (*common.ArchiveInfo, error) {
	arc := &common.ArchiveInfo{
		Platform: Platform,
	}

	tagName := ""
	if len(url.Release) > 0 {
		tagName = url.Release
	}
	if len(url.Tag) > 0 {
		tagName = url.Tag
	}

	if len(tagName) > 0 {
		tag, err := findTag(url.Owner, url.Repo, tagName)
		if err != nil {
			return nil, err
		}
		if tag == nil {
			return nil, errors.New("no matched tag")
		}

		// https://github.com/{owner}/{repo}/archive/refs/tags/{tag}.{format}
		// 这里使用 Archive URL 而不使用 REST API 返回的 URL 可以得到更友好的文件名
		arcUrl := fmt.Sprintf("https://github.com/%s/%s/archive/refs/tags/%s", url.Owner, url.Repo, tagName)

		arc.Name = tagName
		arc.Commit = tag.Commit.SHA
		arc.TarUrl = arcUrl + ".tar.gz"
		arc.ZipUrl = arcUrl + ".zip"
		// arc.TarUrl = tag.TarBallUrl
		// arc.ZipUrl = tag.ZipBallUrl
		return arc, nil
	}

	// releases > tags > branches

	if arc.Commit == "" {
		return nil, errors.New("no commit found")
	}

	return arc, nil
}

func findTag(owner, repo, name string) (*Tag, error) {
	for page := 1; page < 1000; page++ {
		tags, err := GetTags(owner, repo, page)
		if err != nil {
			return nil, err
		}
		if tags == nil || len(tags) <= 0 {
			return nil, errors.New("no tag fetched")
		}
		for _, tag := range tags {
			if tag.Name == name {
				return &tag, nil
			}
		}
	}
	return nil, nil
}
