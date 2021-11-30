package storage

import (
	"context"
	"path"

	"github.com/mises-id/socialsvc/config/env"
)

var (
	storageService IStorageService
	Prefix         = "upload/"
)

func init() {
	switch env.Envs.StorageProvider {
	default:
		storageService = &FileStore{}
	case "local":
		storageService = &FileStore{}
	case "oss":
		storageService = &OSSStorage{}
	}
}

func UploadFile(ctx context.Context, filePath, filename string, file File) error {
	realPath := path.Join(env.Envs.RootPath, Prefix, filePath)
	return storageService.Upload(ctx, realPath, filename, file)
}

type IStorageService interface {
	Upload(ctx context.Context, filePath, filename string, file File) error
}

type File interface {
	Read(p []byte) (n int, err error)
}
