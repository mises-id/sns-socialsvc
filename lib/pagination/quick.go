package pagination

import (
	"github.com/mises-id/socialsvc/lib/db/odm"
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
	Limit  int64  `json:"limit" query:"limit"`
	NextID string `json:"last_id" query:"last_id"`
}

type QuickPaginator struct {
	Limit  int64   `json:"-"`
	NextID string  `json:"-"`
	DB     *odm.DB `json:"-"`
}

func NewQuickPaginator(limit int64, nextID string, db *odm.DB) Paginator {
	if limit == 0 {
		limit = 50
	}

	return &QuickPaginator{
		Limit:  limit,
		NextID: nextID,
		DB:     db,
	}
}

type nextItem struct {
	ID string `bson:"_id,omitempty"`
}

func (p *QuickPaginator) Paginate(dataSource interface{}) (Pagination, error) {
	db := p.DB
	var err error
	if p.NextID != "" {
		hex, err := primitive.ObjectIDFromHex(p.NextID)
		if err != nil {
			return nil, err
		}
		db = db.Where(bson.M{"_id": bson.M{"$lte": hex}})
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
		Limit:  p.Limit,
		NextID: nextID,
	}, nil
}

func (p *QuickPagination) BuildJSONResult() interface{} {
	return p
}

func (p *QuickPagination) GetPageSize() int {
	return int(p.Limit)
}
