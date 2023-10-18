package common

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

type ArchiveInfo struct {
	Platform string
	Name     string
	Commit   string
	TarUrl   string
	ZipUrl   string
}
