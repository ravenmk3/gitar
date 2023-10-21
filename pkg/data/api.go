package data

type DataStore interface {
	Open() error
	Close() error

	RepoExists(repo string) (bool, error)
	SaveRepo(repo string) error

	IsCommitDownloaded(id string) (bool, error)
	SetCommitDownloaded(id string) error

	IsCommitMailed(id string) (bool, error)
	SetCommitMailed(id string) error
}
