package pagination

import (
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PageQuickParams struct {
	Limit  int64  `json:"limit" query:"limit"`
	NextID string `json:"last_id" query:"last_id"`
}

func DefaultQuickParams() *PageQuickParams {
	return &PageQuickParams{
		Limit: 50,
	}
}

func (*PageQuickParams) isPagePrams() {}

func (p *PageQuickParams) GetLimit() int64 {
	if p.Limit <= 0 || p.Limit > 200 {
		return 50
	}
	return p.Limit
}

type QuickPagination struct {
	TotalRecords int64  `json:"total_records" query:"total_records"`
	Limit        int64  `json:"limit" query:"limit"`
	NextID       string `json:"last_id" query:"last_id"`
}

type QuickPaginator struct {
	Limit    int64   `json:"-"`
	NextID   string  `json:"-"`
	SortType string  `json:"-"`
	Count    bool    `json:"-"`
	DB       *odm.DB `json:"-"`
}

type Options func(p *QuickPaginator)

func NewQuickPaginator(limit int64, nextID string, db *odm.DB, opts ...Options) Paginator {
	if limit == 0 {
		limit = 50
	}
	qp := &QuickPaginator{
		Limit:  limit,
		NextID: nextID,
		DB:     db,
	}
	for _, opt := range opts {
		opt(qp)
	}
	return qp
}

func SortAsc() Options {
	return func(p *QuickPaginator) {
		p.SortType = "asc"
	}
}
func IsCount(count bool) Options {
	return func(p *QuickPaginator) {
		p.Count = count
	}
}

type nextItem struct {
	ID string `bson:"_id,omitempty"`
}

func (p *QuickPaginator) Paginate(dataSource interface{}) (Pagination, error) {
	db := p.DB
	var err error
	var count int64
	if p.Count {
		err1 := db.Model(dataSource).Count(&count).Error
		if err1 != nil {
			count = 0
		}
	}
	if p.NextID != "" {
		hex, err := primitive.ObjectIDFromHex(p.NextID)
		if err != nil {
			return nil, err
		}
		db = db.Where(bson.M{"_id": bson.M{"$lte": hex}})
		if p.SortType == "asc" {
			db = db.Where(bson.M{"_id": bson.M{"$gte": hex}})
		}
	}
	err = db.Sort(bson.M{"_id": -1}).Limit(p.Limit).Find(dataSource).Error
	if err != nil {
		return nil, err
	}

	items := make([]*nextItem, 0)
	if err = db.Skip(p.Limit).Limit(1).Find(&items).Error; err != nil {
		return nil, err
	}
	nextID := ""
	if len(items) > 0 {
		nextID = items[0].ID
	}
	return &QuickPagination{
		Limit:        p.Limit,
		NextID:       nextID,
		TotalRecords: count,
	}, nil
}

func (p *QuickPagination) BuildJSONResult() interface{} {
	return p
}

func (p *QuickPagination) GetPageSize() int {
	return int(p.Limit)
}
