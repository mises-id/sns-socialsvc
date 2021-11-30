package storage

import (
	"context"
	"io"
	"os"
	"path"
)

type FileStore struct{}

func (s *FileStore) Upload(ctx context.Context, filePath, filename string, file File) error {
	var err error
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	dst, err := os.Create(path.Join(filePath, filename))
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}
