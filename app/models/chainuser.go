package models

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ChainUser struct {
		ID        primitive.ObjectID   `bson:"_id,omitempty"`
		Misesid   string               `bson:"misesid"`
		Pubkey    string               `bson:"pubkey"`
		Status    enum.ChainUserStatus `bson:"status"`
		TxID      string               `bson:"tx_id"`
		CreatedAt time.Time            `bson:"created_at,omitempty"`
	}
)

func FindChainUser(ctx context.Context, params IAdminParams) (*ChainUser, error) {

	res := &ChainUser{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Last(res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ListChainUser(ctx context.Context, params IAdminParams) ([]*ChainUser, error) {

	res := make([]*ChainUser, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *ChainUser) BeforeCreate(ctx context.Context) error {
	var lc sync.Mutex
	lc.Lock()
	this.ID = primitive.NilObjectID
	this.CreatedAt = time.Now()
	query := db.ODM(ctx).Where(bson.M{"misesid": this.Misesid})

	var c int64
	err := query.Model(this).Count(&c).Error
	lc.Unlock()
	if err != nil {
		return err
	}
	if c > 0 {
		return errors.New("misesid exists")
	}
	return nil
}

func CreateChainUser(ctx context.Context, data *ChainUser) error {

	if err := data.BeforeCreate(ctx); err != nil {
		return err
	}
	_, err := db.DB().Collection("chainusers").InsertOne(ctx, data)

	return err
}

func (m *ChainUser) UpdateTxID(ctx context.Context, tx_id string) error {
	update := bson.M{}
	update["tx_id"] = tx_id
	_, err := db.DB().Collection("chainusers").UpdateOne(ctx, &bson.M{
		"_id": m.ID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}

func (m *ChainUser) UpdateStatus(ctx context.Context, status enum.ChainUserStatus) error {
	update := bson.M{}
	update["status"] = status
	_, err := db.DB().Collection("chainusers").UpdateOne(ctx, &bson.M{
		"_id": m.ID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}
