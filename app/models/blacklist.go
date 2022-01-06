package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Blacklist struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UID        uint64             `bson:"uid,omitempty"`
	TargetUID  uint64             `bson:"target_uid,omitempty"`
	CreatedAt  time.Time          `bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `bson:"updated_at,omitempty"`
	TargetUser *User              `bson:"-"`
}

func ListBlacklist(ctx context.Context, uid uint64, pageParams *pagination.PageQuickParams) ([]*Blacklist, pagination.Pagination, error) {
	if pageParams == nil {
		pageParams = pagination.DefaultQuickParams()
	}
	blacklists := make([]*Blacklist, 0)
	chain := db.ODM(ctx).Where(bson.M{"uid": uid})
	paginator := pagination.NewQuickPaginator(pageParams.Limit, pageParams.NextID, chain)
	page, err := paginator.Paginate(&blacklists)
	if err == nil {
		return nil, nil, err
	}
	return blacklists, page, preloadBlacklistUser(ctx, blacklists...)
}

func FindBlacklist(ctx context.Context, uid, targetUID uint64) (*Blacklist, error) {
	blacklist := &Blacklist{}
	err := db.ODM(ctx).Where(bson.M{"uid": uid, "target_uid": targetUID}).First(blacklist).Error
	if err != nil {
		return nil, err
	}
	return blacklist, preloadBlacklistUser(ctx, blacklist)

}

func CreateBlacklist(ctx context.Context, uid, targetUID uint64) (*Blacklist, error) {
	blacklist, err := FindBlacklist(ctx, uid, targetUID)
	if err == nil {
		return blacklist, nil
	}
	if mongo.ErrNoDocuments != err {
		return nil, err
	}
	blacklist = &Blacklist{
		UID:       uid,
		TargetUID: targetUID,
	}
	err = db.ODM(ctx).Create(blacklist).Error
	if err != nil {
		return nil, err
	}
	return blacklist, preloadBlacklistUser(ctx, blacklist)
}

func DeleteBlacklist(ctx context.Context, uid, targetUID uint64) error {
	blacklist, err := FindBlacklist(ctx, uid, targetUID)
	if err != nil {
		return err
	}
	return db.ODM(ctx).Delete(blacklist, blacklist.ID).Error
}

func preloadBlacklistUser(ctx context.Context, blacklists ...*Blacklist) error {
	userIDs := make([]uint64, len(blacklists))
	for i, blacklist := range blacklists {
		userIDs[i] = blacklist.TargetUID
	}
	users, err := FindUserByIDs(ctx, userIDs...)
	if err != nil {
		return err
	}
	userMap := make(map[uint64]*User)
	for _, user := range users {
		userMap[user.UID] = user
	}
	for _, blacklist := range blacklists {
		blacklist.TargetUser = userMap[blacklist.TargetUID]
	}
	return nil
}
