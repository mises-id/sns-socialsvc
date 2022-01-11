package admin

import (
	"context"
	"errors"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	StatusApi interface {
		PageStatus(ctx context.Context, params *AdminStatusParams) ([]*Status, pagination.Pagination, error)
		ListStatus(ctx context.Context, params *AdminStatusParams) ([]*Status, error)
		FindStatus(ctx context.Context, params *AdminStatusParams) (*Status, error)
		CreateTag(ctx context.Context, id primitive.ObjectID, tag enum.TagType) (*StatusTag, error)
		DeleteTag(ctx context.Context, id primitive.ObjectID, tag enum.TagType) (*StatusTag, error)
	}
	StatusTag struct {
		*models.Status
		*models.Tag
	}
	Status struct {
		*models.Status
	}
	statusApi struct {
	}
)

func NewStatusApi() StatusApi {
	return &statusApi{}
}

func (a *statusApi) FindStatus(ctx context.Context, params *AdminStatusParams) (*Status, error) {
	status, err := models.AdminFindStatus(ctx, params)
	if err != nil {
		return nil, err
	}
	return &Status{status}, nil
}
func (a *statusApi) ListStatus(ctx context.Context, params *AdminStatusParams) ([]*Status, error) {
	statuses, err := models.AdminListStatus(ctx, params)
	if err != nil {
		return nil, err
	}
	result := make([]*Status, len(statuses))
	for i, status := range statuses {
		result[i] = &Status{status}
	}
	return result, nil
}

func (a *statusApi) PageStatus(ctx context.Context, params *AdminStatusParams) ([]*Status, pagination.Pagination, error) {

	statuses, page, err := models.AdminPageStatus(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	result := make([]*Status, len(statuses))
	for i, status := range statuses {
		result[i] = &Status{status}
	}
	return result, page, nil
}

func buildStatusTag(status *models.Status, tag *models.Tag) *StatusTag {
	return &StatusTag{
		Status: status,
		Tag:    tag,
	}
}

func (a *statusApi) CreateTag(ctx context.Context, id primitive.ObjectID, tag enum.TagType) (*StatusTag, error) {

	status := &models.Status{}
	if err := db.ODM(ctx).First(status, bson.M{"_id": id}).Error; err != nil {
		return nil, errors.New("status not found")
	}
	tags := status.Tags
	if index := inArray(tags, tag); index >= 0 {
		return nil, errors.New("tag exists")
	}
	//TODO private status

	switch tag {
	case enum.TagRecommendStatus:
		hide_time := status.HideTime
		if hide_time != nil && time.Now().Unix() > hide_time.Unix() {
			return nil, errors.New("private cannot recommend")
		}
	}

	tags = append(tags, tag)
	status.Tags = tags
	if err := models.UpdateStatusTag(ctx, status); err != nil {
		return nil, err
	}
	params := &models.CreateTagParams{
		TagType:     tag,
		TagableID:   id.String(),
		TagableType: enum.TagableStatus,
	}
	tag_data, err := models.CreateTag(ctx, params)
	if err != nil {
		return nil, err
	}

	return &StatusTag{nil, tag_data}, nil
}

func (a *statusApi) DeleteTag(ctx context.Context, id primitive.ObjectID, tag enum.TagType) (*StatusTag, error) {
	status := &models.Status{}
	if err := db.ODM(ctx).First(status, bson.M{"_id": id}).Error; err != nil {
		return nil, err
	}
	tags := status.Tags
	if index := inArray(tags, tag); index >= 0 {
		tags = append(tags[:index], tags[index+1:]...)
	}
	status.Tags = tags
	if err := models.UpdateStatusTag(ctx, status); err != nil {
		return nil, err
	}
	err := models.DeleteTag(ctx, id.String(), enum.TagableStatus, tag)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func inArray(arr []enum.TagType, item enum.TagType) int {
	for k, v := range arr {
		if v == item {
			return k
		}
	}
	return -1
}
