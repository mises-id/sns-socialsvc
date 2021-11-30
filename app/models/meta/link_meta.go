package meta

type LinkMeta struct {
	Title         string `json:"title"`
	Host          string `json:"host"`
	AttachmentID  uint64 `json:"attachment_id"`
	Link          string `json:"link"`
	AttachmentURL string `json:"-"`
}

func (*LinkMeta) isMetaData() {}
