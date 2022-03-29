package storage

import (
	"context"

	"github.com/mises-id/sns-storagesvc/sdk/service/imgview"
	"github.com/mises-id/sns-storagesvc/sdk/service/imgview/options"
)

var ImageClient IStorage

type imageStorage struct {
	client *imgview.Client
	opts   *options.ImageOptions
}
type (
	ResizeOption struct {
	}
)

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

func (s *imageStorage) ImageResize(opts *ImageOptions) {
	if opts.ResizeOptions != nil {
		s.opts.ResizeOptions = &options.ResizeOptions{
			Resize:     true,
			Width:      opts.ResizeOptions.Width,
			Height:     opts.ResizeOptions.Height,
			ResizeType: opts.ResizeOptions.ResizeType,
		}
	}
}
func (s *imageStorage) ImageCrop(opts *ImageOptions) {
	if opts.CropOptions != nil {
		s.opts.CropOptions = &options.CropOptions{
			Crop:   true,
			Width:  opts.CropOptions.Width,
			Height: opts.CropOptions.Height,
		}
	}
}

func (s *imageStorage) options(ctx context.Context, opts *ImageOptions) {
	s.opts = &options.ImageOptions{}
	if opts.Quality > 0 {
		s.opts.Quality = opts.Quality
	}
	s.ImageResize(opts)
	s.ImageCrop(opts)
}

func (s *imageStorage) GetFileUrlOptions(ctx context.Context, opts *ImageOptions, paths ...string) (map[string]string, error) {
	s.options(ctx, opts)
	result, err := s.client.GetImgUrlList(ctx, &options.ImageViewListInput{
		Path:         paths,
		ImageOptions: s.opts,
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
func (s *imageStorage) GetFileUrl(ctx context.Context, paths ...string) (map[string]string, error) {
	result, err := s.client.GetImgUrlList(ctx, &options.ImageViewListInput{
		Path:         paths,
		ImageOptions: &options.ImageOptions{},
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
