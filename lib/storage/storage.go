package storage

import "context"

type IStorage interface {
	GetFileUrl(ctx context.Context, path ...string) (map[string]string, error)
	GetFileUrlOptions(ctx context.Context, opts *ImageOptions, path ...string) (map[string]string, error)
}

type (
	WatermarkTextOptions struct {
		Watermark bool
		Text      string
		Font      string
		FontSize  int
		Color     string
	}
	CropOptions struct {
		Crop   bool
		Height int
		Width  int
	}
	ResizeOptions struct {
		Resize bool
		//ResizeType
		//fit  resizes the image while keeping aspect ratio to fit given size;
		//fill resizes the image while keeping aspect ratio to fill given size and cropping projecting parts;
		//force resizes the image without keeping aspect ratio;
		ResizeType string
		Height     int
		Width      int
	}
	ImageOptions struct {
		*ResizeOptions
		*CropOptions
		*WatermarkTextOptions
		Format  string //jpeg,png,jpg,webp
		Quality int    //1-100
	}
)
