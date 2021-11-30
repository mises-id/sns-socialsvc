package models

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/mises-id/socialsvc/app/models/enum"
	"github.com/mises-id/socialsvc/config/env"
	"github.com/mises-id/socialsvc/lib/db"
	"github.com/mises-id/socialsvc/lib/storage"
	"go.mongodb.org/mongo-driver/bson"
)

type Attachment struct {
	ID        uint64        `bson:"_id"`
	Filename  string        `bson:"filename,omitempty"`
	FileType  enum.FileType `bson:"file_type"`
	CreatedAt time.Time     `bson:"created_at,omitempty"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty"`
	file      storage.File
}

func (a *Attachment) BeforeCreate(ctx context.Context) error {
	var err error
	a.ID, err = getNextSeq(ctx, "attachmentid")
	if err != nil {
		return err
	}
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Attachment) FileUrl() string {
	return env.Envs.AssetHost + path.Join(storage.Prefix, a.fileFolder(), a.Filename)
}

func (a *Attachment) fileFolder() string {
	if a.ID == 0 {
		return "tmp"
	}
	const attachmentPrefix = "attachment/"
	return path.Join(attachmentPrefix, a.CreatedAt.Format("2006/01/02/"), strconv.Itoa(int(a.ID)))
}

func (a *Attachment) UploadFile(ctx context.Context) error {
	return storage.UploadFile(ctx, a.fileFolder(), a.Filename, a.file)
}

func CreateAttachment(ctx context.Context, tp enum.FileType, filename string, file storage.File) (*Attachment, error) {
	attachment := &Attachment{
		Filename: filename,
		FileType: tp,
		file:     file,
	}
	if err := attachment.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	if err := attachment.UploadFile(ctx); err != nil {
		return nil, err
	}
	_, err := db.DB().Collection("attachments").InsertOne(ctx, attachment)
	return attachment, err
}

func FindAttachmentMap(ctx context.Context, ids []uint64) (map[uint64]*Attachment, error) {
	attachments := make([]*Attachment, 0)
	cursor, err := db.DB().Collection("attachments").Find(ctx,
		bson.M{
			"_id": bson.M{"$in": ids},
		})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &attachments); err != nil {
		return nil, err
	}
	result := make(map[uint64]*Attachment)
	for _, attachment := range attachments {
		result[attachment.ID] = attachment
	}
	return result, nil
}
