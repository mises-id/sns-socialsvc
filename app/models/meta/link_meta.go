package meta

type LinkMeta struct {
	Title     string `bson:"title" json:"title"`
	Host      string `bson:"host" json:"host"`
	Link      string `bson:"link" json:"link"`
	ImagePath string `bson:"image_path" json:"image_path"`
	ImageURL  string `bson:"-"`
}
