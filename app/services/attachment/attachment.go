package attachment

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/mises-id/sns-socialsvc/app/models"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
)

func CreateAttachment(ctx context.Context, fileType, filename string, file multipart.File) (*models.Attachment, error) {
	filenames := strings.Split(filename, "/")
	tp, err := enum.FileTypeFromString(fileType)
	if err != nil {
		return nil, err
	}
	return models.CreateAttachment(ctx, tp, filenames[len(filenames)-1], file)
}
