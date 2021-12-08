package storage

import (
	"context"

	"github.com/mises-id/sns-socialsvc/lib/codes"
)

type S3Storage struct{}

func (s *S3Storage) Upload(ctx context.Context, filePath, filename string, file File) error {
	return codes.ErrUnimplemented
}
