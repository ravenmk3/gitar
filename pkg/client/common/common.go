package common

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

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

type ArchiveResolver interface {
	ResolveArchive(url RepoUrl) (*ArchiveInfo, error)
}

func ResolveArchiveWithRetry(url RepoUrl, resolver ArchiveResolver, maxAttempts int) (*ArchiveInfo, error) {
	for i := 0; i < maxAttempts; i++ {
		info, err := resolver.ResolveArchive(url)
		if err == nil {
			return info, nil
		}
		if errors.Is(err, context.DeadlineExceeded) {
			logrus.Warnf("Error while resolving: %s", err.Error())
			continue
		}
		return nil, err
	}
	return nil, fmt.Errorf("Could not resolve archive after %d attempts", maxAttempts)
}
