package admin

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListStatusParams struct {
	*pagination.TraditionalParams
	StartTime       *time.Time
	EndTime         *time.Time
	MinCommentCount uint64
	MaxCommentCount uint64
	MinLikeCount    uint64
	MaxLikeCount    uint64
	MinForwardCount uint64
	MaxForwardCount uint64
	OrderBy         string
}

type StatusTag struct {
	*models.Status
	*models.Tag
}

type ListStatusTagParams struct {
	*pagination.TraditionalParams
	Tag enum.TagType
}

type StatusApi interface {
	ListTag(ctx context.Context, params *ListStatusTagParams) ([]*StatusTag, pagination.Pagination, error)
	CreateTag(ctx context.Context, id primitive.ObjectID, tag enum.TagType) (*StatusTag, error)
	DeleteTag(ctx context.Context, id primitive.ObjectID, tag enum.TagType) (*StatusTag, error)
}

type statusApi struct {
}

func NewStatusApi() StatusApi {
	return &statusApi{}
}

func (a *statusApi) ListTag(ctx context.Context, params *ListStatusTagParams) ([]*StatusTag, pagination.Pagination, error) {
	return nil, nil, nil
}

func (a *statusApi) CreateTag(ctx context.Context, id primitive.ObjectID, tag enum.TagType) (*StatusTag, error) {
	return nil, nil
}

func (a *statusApi) DeleteTag(ctx context.Context, id primitive.ObjectID, tag enum.TagType) (*StatusTag, error) {
	return nil, nil
}
