package meta

type MetaData struct {
	TextMeta  *TextMeta  `bson:"text_meta"`
	LinkMeta  *LinkMeta  `bson:"link_meta"`
	ImageMeta *ImageMeta `bson:"image_meta"`
}
