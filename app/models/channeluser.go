package models

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ChannelUserExist = errors.New("channel user exist")
)

type (
	ChannelUser struct {
		ID             primitive.ObjectID       `bson:"_id,omitempty"`
		ChannelID      primitive.ObjectID       `bson:"channel_id"`
		ChannelMisesid string                   `bson:"channel_misesid"`
		UID            uint64                   `bson:"uid"`
		ValidState     enum.UserValidState      `bson:"valid_state"`   //default  success failed
		AirdropState   enum.ChannelAirdropState `bson:"airdrop_state"` //default  pending  success failed
		TxID           string                   `bson:"tx_id"`
		Amount         int64                    `bson:"amount"`
		AirdropError   string                   `bson:"airdrop_error"`
		ValidTime      time.Time                `bson:"valid_time"`
		AirdropTime    time.Time                `bson:"airdrop_time"`
		CreatedAt      time.Time                `bson:"created_at"`
		User           *User                    `bson:"-"`
	}
	PageChannelUserInput struct {
		PageParams *pagination.PageQuickParams
		Misesid    string
	}
)

func FindChannelUser(ctx context.Context, params IAdminParams) (*ChannelUser, error) {

	res := &ChannelUser{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Get(res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func FindChannelUserByID(ctx context.Context, id primitive.ObjectID) (*ChannelUser, error) {
	res := &ChannelUser{}
	result := db.DB().Collection("channelusers").FindOne(ctx, &bson.M{
		"_id": id,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	return res, result.Decode(res)
}

func ListChannelUser(ctx context.Context, params IAdminParams) ([]*ChannelUser, error) {

	res := make([]*ChannelUser, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *ChannelUser) BeforeCreate(ctx context.Context) error {
	var lc sync.Mutex
	lc.Lock()
	defer lc.Unlock()
	this.ID = primitive.NilObjectID
	this.CreatedAt = time.Now()
	query := db.ODM(ctx).Where(bson.M{"uid": this.UID})

	var c int64
	err := query.Model(this).Count(&c).Error
	if err != nil {
		return err
	}
	if c > 0 {
		return ChannelUserExist
	}
	return nil
}

func CreateChannelUser(ctx context.Context, data *ChannelUser) (*ChannelUser, error) {

	if err := data.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	res, err := db.DB().Collection("channelusers").InsertOne(ctx, data)
	data.ID = res.InsertedID.(primitive.ObjectID)
	return data, err
}

func (m *ChannelUser) UpdateTxID(ctx context.Context, tx_id string) error {
	update := bson.M{}
	update["tx_id"] = tx_id
	_, err := db.DB().Collection("channelusers").UpdateOne(ctx, &bson.M{
		"_id": m.ID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}
func (m *ChannelUser) UpdateStatusPending(ctx context.Context) error {
	update := bson.M{}
	update["airdrop_state"] = enum.ChannelAirdropPending
	fmt.Println("up id: ", m.ID)
	_, err := db.DB().Collection("channelusers").UpdateOne(ctx, &bson.M{
		"_id": m.ID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}
func (m *ChannelUser) UpdateStatusFailed(ctx context.Context, airdrop_error string) error {
	update := bson.M{}
	update["airdrop_state"] = enum.ChannelAirdropFailed
	update["airdrop_error"] = airdrop_error
	update["airdrop_time"] = time.Now()
	_, err := db.DB().Collection("channelusers").UpdateOne(ctx, &bson.M{
		"_id": m.ID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}
func (m *ChannelUser) UpdateStatusSuccess(ctx context.Context) error {
	update := bson.M{}
	update["airdrop_state"] = enum.ChannelAirdropSuccess
	update["airdrop_time"] = time.Now()
	_, err := db.DB().Collection("channelusers").UpdateOne(ctx, &bson.M{
		"_id": m.ID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}

func (m *ChannelUser) UpdateCreateAirdrop(ctx context.Context, amount int64) error {
	update := bson.M{}
	update["valid_state"] = enum.UserValidSucessed
	update["amount"] = amount
	update["valid_time"] = time.Now()
	_, err := db.DB().Collection("channelusers").UpdateOne(ctx, &bson.M{
		"_id": m.ID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}

func PageChannelUser(ctx context.Context, params *PageChannelUserInput) ([]*ChannelUser, pagination.Pagination, error) {
	if params.PageParams == nil {
		params.PageParams = pagination.DefaultQuickParams()
	}
	res := make([]*ChannelUser, 0)
	chain := db.ODM(ctx)
	and := bson.A{}
	if params.Misesid != "" {
		and = append(and, bson.M{"channel_misesid": utils.AddMisesidProfix(params.Misesid)})
	}
	if len(and) > 0 {
		chain = chain.Where(bson.M{"$and": and})
	}
	paginator := pagination.NewQuickPaginator(params.PageParams.Limit, params.PageParams.NextID, chain, pagination.IsCount(true))
	page, err := paginator.Paginate(&res)
	if err != nil {
		return nil, nil, err
	}
	return res, page, preloadChannelUser(ctx, res...)
}

func preloadChannelUser(ctx context.Context, channel_users ...*ChannelUser) error {
	userIds := make([]uint64, 0)
	for _, channel_user := range channel_users {
		userIds = append(userIds, channel_user.UID)
	}
	users, err := FindUserByIDs(ctx, userIds...)
	if err != nil {
		return err
	}
	userMap := make(map[uint64]*User)
	for _, user := range users {
		userMap[user.UID] = user
	}
	for _, channel_user := range channel_users {
		channel_user.User = userMap[channel_user.UID]
	}
	return nil
}

func CountChannelUser(ctx context.Context, params IAdminParams) (int64, error) {

	var res int64
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Model(&ChannelUser{}).Count(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}
