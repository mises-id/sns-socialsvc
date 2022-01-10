package pagination

import (
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
)

type TraditionalPaginatorAdmin struct {
	PageNum  int64   `json:"-"`
	PageSize int64   `json:"-"`
	Offset   int64   `json:"-"`
	DB       *odm.DB `json:"-"`
}

func NewTraditionalPaginatorAdmin(pageNum, pageSize int64, db *odm.DB) Paginator {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 2 || pageSize > 200 {
		pageSize = 50
	}
	offset := (pageNum - 1) * pageSize

	return &TraditionalPaginatorAdmin{
		PageNum:  pageNum,
		PageSize: pageSize,
		Offset:   offset,
		DB:       db,
	}
}

func (p *TraditionalPaginatorAdmin) Paginate(dataSource interface{}) (Pagination, error) {
	db := p.DB

	var count int64
	err := db.Model(dataSource).Count(&count).Error
	if err != nil {
		return nil, err
	}
	err = db.Limit(p.PageSize).Skip(p.Offset).Find(dataSource).Error
	if err != nil {
		return nil, err
	}
	totalPages := count / p.PageSize
	if count%int64(p.PageSize) > 0 {
		totalPages++
	}

	return &TraditionalPagination{
		TotalRecords: count,
		TotalPages:   totalPages,
		PageSize:     p.PageSize,
		PageNum:      p.PageNum,
	}, nil
}
