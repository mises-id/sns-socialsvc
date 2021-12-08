package pagination

import (
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"go.mongodb.org/mongo-driver/bson"
)

type TraditionalParams struct {
	PageNum  int64 `json:"page_num" query:"page_num"`
	PageSize int64 `json:"page_size" query:"page_size"`
}

func (*TraditionalParams) isPagePrams() {}

type TraditionalPagination struct {
	TotalRecords int64 `json:"total_records" query:"total_records"`
	TotalPages   int64 `json:"total_pages" query:"total_pages"`
	PageNum      int64 `json:"page_num" query:"page_num"`
	PageSize     int64 `json:"page_size" query:"page_size"`
}

type TraditionalPaginator struct {
	PageNum  int64   `json:"-"`
	PageSize int64   `json:"-"`
	Offset   int64   `json:"-"`
	DB       *odm.DB `json:"-"`
}

func DefaultTraditionalParams() *TraditionalParams {
	return &TraditionalParams{
		PageNum:  1,
		PageSize: 50,
	}
}

func NewTraditionalParams(pageNum, pageSize int64) *TraditionalParams {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize <= 2 || pageSize > 200 {
		pageSize = 50
	}
	return &TraditionalParams{
		PageNum:  pageNum,
		PageSize: pageSize,
	}
}

func NewTraditionalPaginator(pageNum, pageSize int64, db *odm.DB) Paginator {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 2 || pageSize > 200 {
		pageSize = 50
	}
	offset := (pageNum - 1) * pageSize

	return &TraditionalPaginator{
		PageNum:  pageNum,
		PageSize: pageSize,
		Offset:   offset,
		DB:       db,
	}
}

func (p *TraditionalPaginator) Paginate(dataSource interface{}) (Pagination, error) {
	db := p.DB

	var count int64
	err := db.Model(dataSource).Count(&count).Error
	if err != nil {
		return nil, err
	}
	err = db.Sort(bson.M{"_id": -1}).Limit(p.PageSize).Skip(p.Offset).Find(dataSource).Error
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

func (p *TraditionalPagination) BuildJSONResult() interface{} {
	return p
}

func (p *TraditionalPagination) GetPageSize() int {
	return int(p.PageSize)
}

func (p *TraditionalPagination) SetPageToken(lastRecordID uint64) {
}
