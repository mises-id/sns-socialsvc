package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Airdrop struct {
		ID        primitive.ObjectID `bson:"_id,omitempty"`
		UID       uint64             `bson:"uid"`
		Misesid   string             `bson:"misesid,omitempty"`
		Type      enum.AirdropType   `bson:"type"`
		Coin      int64              `bson:"coin"`
		TxID      string             `bson:"tx_id"`
		Status    enum.AirdropStatus `bson:"status"`
		FinishAt  time.Time          `bson:"finish_at"`
		CreatedAt time.Time          `bson:"created_at"`
	}
)

func FindAirdrop(ctx context.Context, params IAdminParams) (*Airdrop, error) {

	res := &Airdrop{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Get(res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ListAirdrop(ctx context.Context, params IAdminParams) ([]*Airdrop, error) {

	res := make([]*Airdrop, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func CreateAirdropMany(ctx context.Context, data []*Airdrop) error {
	if len(data) == 0 {
		return nil
	}
	var in []interface{}
	for _, v := range data {
		in = append(in, v)
	}
	_, err := db.DB().Collection("airdrops").InsertMany(ctx, in)

	return err
}

func (m *Airdrop) UpdateTxID(ctx context.Context, tx_id string) error {
	update := bson.M{}
	update["tx_id"] = tx_id
	update["status"] = enum.AirdropPending
	_, err := db.DB().Collection("airdrops").UpdateOne(ctx, &bson.M{
		"_id": m.ID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}
func (m *Airdrop) UpdateStatusPending(ctx context.Context) error {
	update := bson.M{}
	update["status"] = enum.AirdropPending
	_, err := db.DB().Collection("airdrops").UpdateOne(ctx, &bson.M{
		"_id": m.ID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}

func (m *Airdrop) UpdateStatus(ctx context.Context, status enum.AirdropStatus) error {
	update := bson.M{}
	update["status"] = status
	update["finish_at"] = time.Now()
	_, err := db.DB().Collection("airdrops").UpdateOne(ctx, &bson.M{
		"_id": m.ID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}

func (m *Airdrop) Update(ctx context.Context) error {

	update := bson.M{}
	if m.TxID != "" {
		update["tx_id"] = m.TxID
	}
	if m.Status != 0 {
		update["status"] = m.Status
		update["finish_at"] = time.Now()
	}
	_, err := db.DB().Collection("airdrops").UpdateOne(ctx, &bson.M{
		"_id": m.ID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err

}
