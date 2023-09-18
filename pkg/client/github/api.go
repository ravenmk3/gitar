package github

import (
	"fmt"

	"gitar/pkg/utils"
)

const ApiBaseUrl = "https://api.github.com"

type TagCommit struct {
	SHA string `json:"sha,omitempty"`
	Url string `json:"url,omitempty"`
}

type Tag struct {
	Name       string    `json:"name,omitempty"`
	ZipBallUrl string    `json:"zipball_url,omitempty"`
	TarBallUrl string    `json:"tarball_url,omitempty"`
	Commit     TagCommit `json:"commit,omitempty"`
	NodeId     string    `json:"node_id,omitempty"`
}

func GetTags(owner string, repo string, page int) ([]Tag, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/tags?page=%d", ApiBaseUrl, owner, repo, page)
	tags := []Tag{}
	err := utils.HttpGetJson(url, &tags)
	return tags, err
}
