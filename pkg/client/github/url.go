package github

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gitar/pkg/client/common"
	"gitar/pkg/utils"
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
		info.Branch = match[3]
		info.RefName = match[3]
		return info, nil
	}

	return nil, fmt.Errorf("unsupported GitHub url %s", url)
}

func ResolveGithubArchive(url common.RepoUrl) (*common.ArchiveInfo, error) {
	tagName := ""
	if len(url.Release) > 0 {
		tagName = url.Release
	}
	if len(url.Tag) > 0 {
		tagName = url.Tag
	}

	if len(tagName) > 0 {
		return resolveArchiveByTag(url, tagName)
	}
	if len(url.Branch) > 0 {
		return resolveArchiveByBranch(url)
	}
	if len(url.Commit) > 0 {
		return resolveArchiveByCommit(url)
	}

	release, err := findAnyRelease(url.Owner, url.Repo)
	if err != nil {
		return nil, err
	}
	if release != nil {
		return resolveArchiveByTag(url, release.TagName)
	}

	branch, err := findAnyWellKnownBranch(url.Owner, url.Repo)
	if err != nil {
		return nil, err
	}
	if branch != nil {
		return branchToArchive(url, *branch)
	}

	return nil, errors.New("could not resolve archive")
}

func branchToArchive(url common.RepoUrl, branch Branch) (*common.ArchiveInfo, error) {
	arc := &common.ArchiveInfo{
		Platform: Platform,
	}

	// https://github.com/{owner}/{repo}/archive/{commit-sha}.{format}
	// 使用 Commit ID 保证在下载时和 API 查到的保持一致

	arcUrl := fmt.Sprintf("https://github.com/%s/%s/archive/%s", url.Owner, url.Repo, branch.Commit.SHA)

	arc.Name = fmt.Sprintf("%s-%s-%s", url.Repo, branch.Name, branch.Commit.SHA[:7])
	arc.Commit = branch.Commit.SHA
	arc.TarUrl = arcUrl + ".tar.gz"
	arc.ZipUrl = arcUrl + ".zip"

	return validateArchive(arc)
}

func validateArchive(arc *common.ArchiveInfo) (*common.ArchiveInfo, error) {
	if arc.Commit == "" {
		return nil, errors.New("no commit found")
	}
	return arc, nil
}

func resolveArchiveByCommit(url common.RepoUrl) (*common.ArchiveInfo, error) {
	arc := &common.ArchiveInfo{
		Platform: Platform,
	}

	arcUrl := fmt.Sprintf("https://github.com/%s/%s/archive/%s", url.Owner, url.Repo, url.Commit)

	arc.Name = fmt.Sprintf("%s-%s", url.Repo, url.Commit[:7])
	arc.Commit = url.Commit
	arc.TarUrl = arcUrl + ".tar.gz"
	arc.ZipUrl = arcUrl + ".zip"

	return validateArchive(arc)
}

func resolveArchiveByBranch(url common.RepoUrl) (*common.ArchiveInfo, error) {
	arc := &common.ArchiveInfo{
		Platform: Platform,
	}

	branch, err := findBranch(url.Owner, url.Repo, url.Branch)
	if err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, errors.New("no matched branch")
	}

	// https://github.com/{owner}/{repo}/archive/{commit-sha}.{format}
	// 使用 Commit ID 保证在下载时和 API 查到的保持一致

	arcUrl := fmt.Sprintf("https://github.com/%s/%s/archive/%s", url.Owner, url.Repo, branch.Commit.SHA)

	arc.Name = fmt.Sprintf("%s-%s-%s", url.Repo, branch.Name, branch.Commit.SHA[:7])
	arc.Commit = branch.Commit.SHA
	arc.TarUrl = arcUrl + ".tar.gz"
	arc.ZipUrl = arcUrl + ".zip"

	return validateArchive(arc)
}

func resolveArchiveByTag(url common.RepoUrl, tagName string) (*common.ArchiveInfo, error) {
	arc := &common.ArchiveInfo{
		Platform: Platform,
	}

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
	arcName := tagName
	if !strings.HasPrefix(tagName, url.Repo) {
		arcName = fmt.Sprintf("%s-%s", url.Repo, tagName)
	}

	arc.Name = arcName
	arc.Commit = tag.Commit.SHA
	arc.TarUrl = arcUrl + ".tar.gz"
	arc.ZipUrl = arcUrl + ".zip"

	return validateArchive(arc)
}

func findAnyWellKnownBranch(owner, repo string) (*Branch, error) {
	desired := utils.NewStringSet([]string{"master", "main", "trunk", "release", "develop"})

	for page := 1; page < 1000; page++ {
		items, err := GetBranches(owner, repo, page)
		if err != nil {
			return nil, err
		}
		if items == nil || len(items) <= 0 {
			return nil, errors.New("no branch fetched")
		}
		for _, item := range items {
			if desired.Contains(item.Name) {
				return &item, nil
			}
		}
	}

	return nil, errors.New("no matched branch")
}

func findAnyRelease(owner, repo string) (*Release, error) {
	items, err := GetReleases(owner, repo, 1)
	if err != nil {
		return nil, err
	}
	if items == nil || len(items) <= 0 {
		return nil, nil
	}
	return &items[0], nil
}

func findBranch(owner, repo, name string) (*Branch, error) {
	for page := 1; page < 1000; page++ {
		items, err := GetBranches(owner, repo, page)
		if err != nil {
			return nil, err
		}
		if items == nil || len(items) <= 0 {
			return nil, errors.New("no branch fetched")
		}
		for _, item := range items {
			if item.Name == name {
				return &item, nil
			}
		}
	}
	return nil, errors.New("no matched branch")
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

	return nil, errors.New("no matched tag")
}
