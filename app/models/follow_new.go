package models

import (
	"context"

	"github.com/mises-id/sns-socialsvc/lib/db"
)

type ()

func NewListFollow(ctx context.Context, params IAdminParams) ([]*Follow, error) {
	res := make([]*Follow, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func CountFollow(ctx context.Context, params IAdminParams) (int64, error) {

	var res int64
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Model(&Follow{}).Count(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}
