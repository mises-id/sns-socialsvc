package recommend

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListStatus(ctx context.Context, uid uint64, num int) ([]primitive.ObjectID, error) {

	return listStatus(ctx, uid, num)

}

func listStatus(ctx context.Context, uid uint64, num int) ([]primitive.ObjectID, error) {

	//find recommend status
	return nil, nil

}
