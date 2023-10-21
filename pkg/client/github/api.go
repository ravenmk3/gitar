package github

import (
	"fmt"
	"time"

	"gitar/pkg/utils"
)

const ApiBaseUrl = "https://api.github.com"

type CommitLite struct {
	SHA string `json:"sha,omitempty"`
	Url string `json:"url,omitempty"`
}

type Tag struct {
	Name       string     `json:"name,omitempty"`
	ZipBallUrl string     `json:"zipball_url,omitempty"`
	TarBallUrl string     `json:"tarball_url,omitempty"`
	Commit     CommitLite `json:"commit,omitempty"`
	NodeId     string     `json:"node_id,omitempty"`
}

type Branch struct {
	Name      string     `json:"name,omitempty"`
	Commit    CommitLite `json:"commit"`
	Protected bool       `json:"protected,omitempty"`
}

type Release struct {
	Id          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	TagName     string    `json:"tag_name,omitempty"`
	Draft       bool      `json:"draft,omitempty"`
	PreRelease  bool      `json:"prerelease,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	ZipBallUrl  string    `json:"zipball_url,omitempty"`
	TarBallUrl  string    `json:"tarball_url,omitempty"`
}

func GetTags(owner string, repo string, page int) ([]Tag, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/tags?page=%d", ApiBaseUrl, owner, repo, page)
	items := []Tag{}
	err := utils.HttpGetJson(url, &items)
	return items, err
}

func GetBranches(owner string, repo string, page int) ([]Branch, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/branches?page=%d", ApiBaseUrl, owner, repo, page)
	items := []Branch{}
	err := utils.HttpGetJson(url, &items)
	return items, err
}

func GetReleases(owner string, repo string, page int) ([]Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases?page=%d", ApiBaseUrl, owner, repo, page)
	items := []Release{}
	err := utils.HttpGetJson(url, &items)
	return items, err
}
