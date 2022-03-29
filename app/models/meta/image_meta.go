package meta

type ImageMeta struct {
	Images         []string `bson:"images"`
	ImageURLs      []string `bson:"-"`
	ThumbImageURLs []string `bson:"-"`
}
