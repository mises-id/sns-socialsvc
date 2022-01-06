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
	//ListTag(ctx context.Context, params *ListUserTagParams) ([]*UserTag, pagination.Pagination, error)
	CreateTag(ctx context.Context, uid uint64, in *CreateUserTagInput) (*UserTag, error)
	DeleteTag(ctx context.Context, uid uint64, tagArr ...enum.TagType) (*UserTag, error)
}

type userApi struct {
}

func NewUserApi() UserApi {
	return &userApi{}
}

func (a *userApi) FindUser(ctx context.Context, params *AdminUserParams) (*User, error) {
	user, err := models.AdminFindUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &User{user}, nil
}
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

	if params.Tag == enum.TagBlank {
		return a.listAllUser(ctx, params.TraditionalParams)
	}
	tags, page, err := models.PageTag(ctx, &models.PageTagParams{
		PageParams: params.TraditionalParams,
		TagParams: models.TagParams{
			TagableType: enum.TagableUser,
			TagType:     params.Tag,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	userIDs := make([]uint64, len(tags))
	for i, tag := range tags {
		uid, _ := strconv.Atoi(tag.TagableID)
		userIDs[i] = uint64(uid)
	}
	users, err := models.ListUserByIDs(ctx, userIDs...)
	if err != nil {
		return nil, nil, err
	}
	userMap := make(map[uint64]*models.User)
	for _, user := range users {
		userMap[user.UID] = user
	}
	result := make([]*UserTag, len(users))
	for i, tag := range tags {
		uid, _ := strconv.Atoi(tag.TagableID)
		result[i] = buildUserTag(userMap[uint64(uid)], tag)
	}
	return result, page, nil
}
*/
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
		TagableID:   string(rune(uid)),
		TagableType: enum.TagableUser,
		Note:        in.Note,
	}
	tag_data, err := models.CreateTag(ctx, params)
	if err != nil {
		return nil, err
	}
	//delete
	if err := models.DeleteTagsByTagtypes(ctx, string(rune(uid)), enum.TagableUser, deleteTags...); err != nil {
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
	err := models.DeleteTagsByTagtypes(ctx, string(rune(uid)), enum.TagableUser, tagArr...)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (a *userApi) listAllUser(ctx context.Context, params *pagination.TraditionalParams) ([]*UserTag, pagination.Pagination, error) {
	users := make([]*models.User, 0)
	chain := db.ODM(ctx)
	paginator := pagination.NewTraditionalPaginator(params.PageNum, params.PageSize, chain)
	page, err := paginator.Paginate(&users)
	if err != nil {
		return nil, nil, err
	}
	userIDs := make([]string, len(users))
	for i, user := range users {
		userIDs[i] = strconv.Itoa(int(user.UID))
	}
	tags, err := models.ListTagByTagables(ctx, enum.TagableUser, userIDs...)
	if err != nil {
		return nil, nil, err
	}
	tagMap := make(map[uint64]*models.Tag)
	for _, tag := range tags {
		uid, _ := strconv.Atoi(tag.TagableID)
		tagMap[uint64(uid)] = tag
	}
	result := make([]*UserTag, len(users))
	for i, user := range users {
		result[i] = buildUserTag(user, tagMap[user.UID])
	}
	return result, page, nil
}

func buildUserTag(user *models.User, tag *models.Tag) *UserTag {
	return &UserTag{
		User: user,
		Tag:  tag,
	}
}
