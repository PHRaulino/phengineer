package discovery

import "context"

type RepoDiscovery interface {
	ListFiles(ctx context.Context) []DicoveredFiles
}
