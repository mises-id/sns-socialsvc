package storage

import (
	"context"

	"github.com/mises-id/sns-storagesvc/sdk/service/imgview"
	"github.com/mises-id/sns-storagesvc/sdk/service/imgview/options"
)

var ImageClient IStorage

type imageStorage struct {
	client *imgview.Client
}

func SetupImageStorage(host, key, salt string) {
	ImageClient = &imageStorage{
		client: imgview.New(
			imgview.Options{
				Key:  key,
				Salt: salt,
				Host: host,
			},
		),
	}
}

func (s *imageStorage) GetFileUrl(ctx context.Context, paths ...string) (map[string]string, error) {
	result, err := s.client.GetImgUrlList(ctx, &options.ImageViewListInput{
		Path: paths,
	})
	if err != nil {
		return nil, err
	}
	imageMap := make(map[string]string)
	for i, url := range result.Url {
		imageMap[paths[i]] = url
	}
	return imageMap, nil
}
