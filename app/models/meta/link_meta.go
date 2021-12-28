package meta

type LinkMeta struct {
	Title     string `bson:"title"`
	Host      string `bson:"host"`
	Link      string `bson:"link"`
	ImagePath string `bson:"image_path"`
	ImageURL  string `bson:"-"`
}
