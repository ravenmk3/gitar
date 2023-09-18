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

type ArchiveUrl struct {
	Platform string
	Commit   string
	TarUrl   string
	ZipUrl   string
}
