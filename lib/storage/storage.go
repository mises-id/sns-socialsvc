package storage

import "context"

type IStorage interface {
	GetFileUrl(ctx context.Context, path ...string) (map[string]string, error)
}
