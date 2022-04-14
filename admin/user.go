package admin

import (
	"context"
	"errors"
	"strconv"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	UserTag struct {
		*models.User
		*models.Tag
	}
	User struct {
		*models.User
	}
	CreateUserTagInput struct {
		Tag  enum.TagType
		Note string
	}
)
type UserApi interface {
	PageUser(ctx context.Context, params *AdminUserParams) ([]*User, pagination.Pagination, error)
	ListUser(ctx context.Context, params *AdminUserParams) ([]*User, error)
	FindUser(ctx context.Context, params *AdminUserParams) (*User, error)
	//user tag
	ListTag(ctx context.Context, params *AdminTagParams) ([]*UserTag, error)
	FindTag(ctx context.Context, params *AdminTagParams) (*UserTag, error)
	CreateTag(ctx context.Context, uid uint64, in *CreateUserTagInput) (*UserTag, error)
	DeleteTag(ctx context.Context, uid uint64, tagArr ...enum.TagType) (*UserTag, error)
}

type userApi struct {
}

func NewUserApi() UserApi {
	return &userApi{}
}

//get one user
func (a *userApi) FindUser(ctx context.Context, params *AdminUserParams) (*User, error) {
	user, err := models.AdminFindUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &User{user}, nil
}

//get list user
func (a *userApi) ListUser(ctx context.Context, params *AdminUserParams) ([]*User, error) {
	users, err := models.AdminListUser(ctx, params)
	if err != nil {
		return nil, err
	}
	result := make([]*User, len(users))
	for i, user := range users {
		result[i] = &User{user}
	}
	return result, nil
}

//get page user
func (a *userApi) PageUser(ctx context.Context, params *AdminUserParams) ([]*User, pagination.Pagination, error) {

	users, page, err := models.AdminPageUser(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	result := make([]*User, len(users))
	for i, user := range users {
		result[i] = &User{user}
	}
	return result, page, nil
}

//create tag
func (a *userApi) CreateTag(ctx context.Context, uid uint64, in *CreateUserTagInput) (*UserTag, error) {

	tag := in.Tag

	user := &models.User{}
	if err := db.ODM(ctx).First(user, bson.M{"_id": uid}).Error; err != nil {
		return nil, errors.New("user not found")
	}
	tags := user.Tags
	if index := inArray(tags, tag); index >= 0 {
		return nil, errors.New("tag exists")
	}
	tags = append(tags, tag)
	//if tag in ['star_user','problem_user']
	onlyTagArr := []enum.TagType{enum.TagStarUser, enum.TagProblemUser}
	var deleteTags []enum.TagType
	if onlyIndex := inArray(onlyTagArr, tag); onlyIndex >= 0 {
		deleteTags = append(onlyTagArr[:onlyIndex], onlyTagArr[onlyIndex+1:]...)
	}
	//delete Tags
	for _, tag := range deleteTags {
		if index := inArray(tags, tag); index >= 0 {
			tags = append(tags[:index], tags[index+1:]...)
		}
	}
	user.Tags = tags
	if err := models.UpdateUserTag(ctx, user); err != nil {
		return nil, err
	}
	params := &models.CreateTagParams{
		TagType:     tag,
		TagableID:   strconv.Itoa(int(uid)),
		TagableType: enum.TagableUser,
		Note:        in.Note,
	}
	tag_data, err := models.CreateTag(ctx, params)
	if err != nil {
		return nil, err
	}
	//delete
	if err := models.DeleteTagsByTagtypes(ctx, strconv.Itoa(int(uid)), enum.TagableUser, deleteTags...); err != nil {
		return nil, err
	}
	return &UserTag{nil, tag_data}, nil

}

//delete tags
func (a *userApi) DeleteTag(ctx context.Context, uid uint64, tagArr ...enum.TagType) (*UserTag, error) {
	if tagArr == nil || len(tagArr) == 0 {
		return nil, nil
	}
	user := &models.User{}
	if err := db.ODM(ctx).First(user, bson.M{"_id": uid}).Error; err != nil {
		return nil, err
	}
	tags := user.Tags
	for _, tag := range tagArr {
		if index := inArray(tags, tag); index >= 0 {
			tags = append(tags[:index], tags[index+1:]...)
		}
	}
	user.Tags = tags
	if err := models.UpdateUserTag(ctx, user); err != nil {
		return nil, err
	}
	err := models.DeleteTagsByTagtypes(ctx, strconv.Itoa(int(uid)), enum.TagableUser, tagArr...)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//get user tags
func (a *userApi) ListTag(ctx context.Context, params *AdminTagParams) ([]*UserTag, error) {

	params.TagableType = enum.TagableUser
	tags, err := models.AdminListTag(ctx, params)
	if err != nil {
		return nil, err
	}
	result := make([]*UserTag, len(tags))
	for i, tag := range tags {
		result[i] = &UserTag{Tag: tag}
	}
	return result, nil

}

//find one tag
func (a *userApi) FindTag(ctx context.Context, params *AdminTagParams) (*UserTag, error) {

	params.TagableType = enum.TagableUser
	tag, err := models.AdminFindTag(ctx, params)
	if err != nil {
		return nil, err
	}

	return &UserTag{Tag: tag}, nil

}
