package message

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeCommentMeta struct {
	UID             uint64             `bson:"uid,omitempty"`
	CommentID       primitive.ObjectID `bson:"comment_id,omitempty"`
	CommentUsername string             `bson:"comment_username,omitempty"`
	CommentContent  string             `bson:"comment_content,omitempty"`
}
