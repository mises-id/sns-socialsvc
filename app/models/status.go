package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/message"
	"github.com/mises-id/sns-socialsvc/app/models/meta"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/storage"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status struct {
	meta.MetaData
	ID                    primitive.ObjectID `bson:"_id,omitempty"`
	ParentID              primitive.ObjectID `bson:"parent_id,omitempty"`
	OriginID              primitive.ObjectID `bson:"origin_id,omitempty"`
	UID                   uint64             `bson:"uid,omitempty"`
	FromType              enum.FromType      `bson:"from_type"`
	StatusType            enum.StatusType    `bson:"status_type"`
	Content               string             `bson:"content,omitempty" validate:"min=0,max=4000"`
	CommentsCount         uint64             `bson:"comments_count,omitempty"`
	LikesCount            uint64             `bson:"likes_count,omitempty"`
	ForwardsCount         uint64             `bson:"forwards_count,omitempty"`
	HideTime              *time.Time         `bson:"hide_time"`
	Score                 int64              `bson:"score"`
	Tags                  []enum.TagType     `bson:"tags"`
	DeletedAt             *time.Time         `bson:"deleted_at,omitempty"`
	CreatedAt             time.Time          `bson:"created_at,omitempty"`
	UpdatedAt             time.Time          `bson:"updated_at,omitempty"`
	User                  *User              `bson:"-"`
	IsLiked               bool               `bson:"-"`
	IsPublic              bool               `bson:"-"`
	ParentStatusIsDeleted bool               `bson:"-"`
	ParentStatusIsBlacked bool               `bson:"-"`
	ParentStatus          *Status            `bson:"-"`
	OriginStatus          *Status            `bson:"-"`
}

func (s *Status) validate(ctx context.Context) error {
	err := Validate.Struct(s)
	if err != nil {
		return codes.ErrUnprocessableEntity
	}
	return nil
}

func (s *Status) BeforeCreate(ctx context.Context) error {
	s.CreatedAt = time.Now()
	s.Score = time.Now().UnixMilli()
	s.UpdatedAt = time.Now()
	var err error
	if !s.ParentID.IsZero() {
		s.ParentStatus, err = FindStatus(ctx, s.ParentID)
		if err != nil {
			return err
		}
		s.OriginID = s.ParentStatus.OriginID
		if s.OriginID.IsZero() {
			s.OriginID = s.ParentID
		}
	}
	if !s.OriginID.IsZero() {
		s.OriginStatus, err = FindStatus(ctx, s.OriginID)
		if err != nil {
			return err
		}
	}
	return s.validate(ctx)
}

func (s *Status) ContentSummary() string {
	content := []rune(s.Content)
	if len(content) < 40 {
		return string(content)
	}
	return string(content[:40])
}

func (s *Status) FirstImage() string {
	if s.ImageMeta != nil && len(s.ImageMeta.Images) > 0 {
		return s.ImageMeta.Images[0]
	}
	return ""
}

func (s *Status) AfterCreate(ctx context.Context) error {
	var err error
	counterKey := s.FromType.CounterKey()
	if s.ParentStatus != nil {
		err = s.ParentStatus.IncStatusCounter(ctx, counterKey)
		if err != nil {
			return err
		}
		//forward public status  notify parent status user
		if s.HideTime == nil {
			message := &CreateMessageParams{
				UID:         s.ParentStatus.UID,
				FromUID:     s.UID,
				StatusID:    s.ID,
				MessageType: enum.NewForward,
				MetaData: &message.MetaData{
					ForwardMeta: &message.ForwardMeta{
						UID:            s.UID,
						StatusID:       s.ID,
						ForwardContent: s.Content,
						ContentSummary: s.ParentStatus.ContentSummary(),
						ImagePath:      s.ParentStatus.FirstImage(),
					},
				},
			}
			_, err = CreateMessage(ctx, message)
			if err != nil {
				return err
			}
		}
	}
	//public status notify fans
	if s.HideTime == nil {
		return s.notifyFans(ctx)
	}
	return nil

}

func (s *Status) GetHideTime() uint64 {
	if s.HideTime == nil {
		return 0
	}
	return uint64(s.HideTime.Unix())
}

func (s *Status) UpdateHideTime(ctx context.Context, isPrivate bool, showDuration int64) error {
	var t *time.Time
	if isPrivate {
		hideTime := time.Now().Add(time.Second * time.Duration(showDuration))
		t = &hideTime
	}
	_, err := db.DB().Collection("statuses").UpdateMany(ctx, bson.M{"_id": s.ID},
		bson.D{{
			Key: "$set",
			Value: bson.D{{
				Key:   "hide_time",
				Value: t,
			}}},
		})
	return err
}

func (s *Status) notifyFans(ctx context.Context) error {
	t := time.Now()
	_, err := db.DB().Collection("follows").UpdateMany(ctx, bson.M{"to_uid": s.UID},
		bson.D{{
			Key: "$set",
			Value: bson.D{{
				Key:   "is_read",
				Value: false,
			}, {
				Key:   "latest_post_time",
				Value: &t,
			}, {
				Key:   "updated_at",
				Value: t,
			}}},
		})
	return err
}

func (s *Status) IncStatusCounter(ctx context.Context, counterKey string, values ...int) error {
	if counterKey == "" {
		return nil
	}
	value := 1
	if len(values) > 0 {
		value = values[0]
	}
	if value < 0 {
		switch counterKey {
		case "likes_count":
			if int(s.LikesCount)+value < 0 {
				value = -int(s.LikesCount)
			}
		case "comments_count":
			if int(s.CommentsCount)+value < 0 {
				value = -int(s.CommentsCount)
			}
		}
	}
	return db.DB().Collection("statuses").FindOneAndUpdate(ctx, bson.M{"_id": s.ID},
		bson.D{{
			Key: "$inc",
			Value: bson.D{{
				Key:   counterKey,
				Value: value,
			}}},
		}).Err()
}

func FindStatus(ctx context.Context, id primitive.ObjectID) (*Status, error) {
	status := &Status{}
	err := db.ODM(ctx).First(status, bson.M{"_id": id}).Error
	if err != nil {
		return nil, err
	}
	return status, PreloadStatusData(ctx, true, status)
}

func FindStatusData(ctx context.Context, id primitive.ObjectID) (*Status, error) {
	status := &Status{}
	err := db.ODM(ctx).First(status, bson.M{"_id": id}).Error
	if err != nil {
		return nil, err
	}
	return status, nil
}

func PreloadStatusData(ctx context.Context, loadRelated bool, statuses ...*Status) error {
	err := preloadAttachment(ctx, statuses...)
	if err != nil {
		return err
	}
	if err = preloadImage(ctx, statuses...); err != nil {
		return err
	}
	if err = preloadStatusUser(ctx, statuses...); err != nil {
		return err
	}
	if err = preloadStatusLikeState(ctx, statuses...); err != nil {
		return err
	}
	if loadRelated {
		err = preloadRelatedStatus(ctx, statuses...)
		if err != nil {
			return err
		}
	}
	return nil
}

func preloadStatusLikeState(ctx context.Context, statuses ...*Status) error {
	currentUID, ok := ctx.Value(utils.CurrentUIDKey{}).(uint64)
	if !ok || currentUID == 0 {
		return nil
	}
	statusIDs := make([]primitive.ObjectID, len(statuses))
	for i, status := range statuses {
		statusIDs[i] = status.ID
	}
	likeMap, err := GetLikeMap(ctx, currentUID, statusIDs, enum.LikeStatus, false)
	if err != nil {
		return err
	}
	for _, status := range statuses {
		status.IsLiked = likeMap[status.ID] != nil
	}
	return nil
}

type CreateStatusParams struct {
	UID          uint64
	ParentID     primitive.ObjectID
	StatusType   enum.StatusType
	FromType     enum.FromType
	Content      string
	IsPrivate    bool
	ShowDuration int64
	MetaData     meta.MetaData
}

func CreateStatus(ctx context.Context, params *CreateStatusParams) (*Status, error) {
	status := &Status{
		UID:        params.UID,
		StatusType: params.StatusType,
		FromType:   params.FromType,
		ParentID:   params.ParentID,
		Content:    params.Content,
		MetaData:   params.MetaData,
	}
	var err error
	if err = status.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	if params.IsPrivate {
		t := status.CreatedAt.Add(time.Duration(params.ShowDuration-10) * time.Second)
		status.HideTime = &t
	}
	if err = db.ODM(ctx).Create(status).Error; err != nil {
		return nil, err
	}
	if err = status.AfterCreate(ctx); err != nil {
		return nil, err
	}
	if err = preloadRelatedStatus(ctx, status); err != nil {
		return nil, err
	}
	if err = preloadAttachment(ctx, status); err != nil {
		return nil, err
	}
	if err = preloadImage(ctx, status); err != nil {
		return nil, err
	}
	return status, preloadStatusUser(ctx, status)
}

func DeleteStatus(ctx context.Context, id primitive.ObjectID) error {
	_, err := db.DB().Collection("statuses").DeleteOne(ctx, bson.M{"_id": id})
	return err
}

type ListStatusParams struct {
	UIDs           []uint64
	CurrentUID     uint64
	ParentStatusID primitive.ObjectID
	FromTypes      []enum.FromType
	OnlyShow       bool
	PageParams     *pagination.PageQuickParams
}

func ListStatus(ctx context.Context, params *ListStatusParams) ([]*Status, pagination.Pagination, error) {
	if params.PageParams == nil {
		params.PageParams = pagination.DefaultQuickParams()
	}
	statuses := make([]*Status, 0)
	chain := db.ODM(ctx)
	and := bson.A{}
	if params.UIDs != nil && len(params.UIDs) > 0 {
		//chain = chain.Where(bson.M{"uid": bson.M{"$in": params.UIDs}})
		and = append(and, bson.M{"uid": bson.M{"$in": params.UIDs}})
	}
	if !params.ParentStatusID.IsZero() {
		chain = chain.Where(bson.M{"parent_id": params.ParentStatusID})
	}
	if params.FromTypes != nil {
		chain = chain.Where(bson.M{"from_type": bson.M{"$in": params.FromTypes}})
	}
	if params.OnlyShow {
		orFilter := bson.A{bson.M{"hide_time": nil}, bson.M{"hide_time": bson.M{"$gt": time.Now().UTC()}}}
		if params.CurrentUID > 0 {
			orFilter = append(orFilter, bson.M{"uid": params.CurrentUID})
		}
		chain = chain.Where(bson.M{"$or": orFilter})
	}
	blockedUIDs, err := ListBlockedUIDs(ctx)
	if err != nil {
		return nil, nil, err
	}
	if len(blockedUIDs) > 0 {
		//chain = chain.Where(bson.M{"uid": bson.M{"$nin": blockedUIDs}})
		and = append(and, bson.M{"uid": bson.M{"$nin": blockedUIDs}})
	}
	if len(and) > 0 {
		chain = chain.Where(bson.M{"$and": and})
	}
	paginator := pagination.NewQuickPaginator(params.PageParams.Limit, params.PageParams.NextID, chain)
	page, err := paginator.Paginate(&statuses)
	if err != nil {
		return nil, nil, err
	}
	return statuses, page, PreloadStatusData(ctx, true, statuses...)
}

func FindStatusByIDs(ctx context.Context, ids ...primitive.ObjectID) ([]*Status, error) {
	statuses := make([]*Status, 0)
	err := db.ODM(ctx).Where(bson.M{"_id": bson.M{"$in": ids}}).Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, PreloadStatusData(ctx, true, statuses...)
}

func ListCommentStatus(ctx context.Context, statusID primitive.ObjectID, pageParams *pagination.PageQuickParams) ([]*Status, pagination.Pagination, error) {
	if pageParams == nil {
		pageParams = pagination.DefaultQuickParams()
	}
	statuses := make([]*Status, 0)
	chain := db.ODM(ctx).Where(bson.M{"parent_id": statusID, "from_type": enum.FromComment})
	paginator := pagination.NewQuickPaginator(pageParams.Limit, pageParams.NextID, chain)
	page, err := paginator.Paginate(&statuses)
	if err != nil {
		return nil, nil, err
	}
	return statuses, page, PreloadStatusData(ctx, true, statuses...)
}

func preloadStatusUser(ctx context.Context, statuses ...*Status) error {
	userIds := make([]uint64, 0)
	for _, status := range statuses {
		userIds = append(userIds, status.UID)
	}
	users, err := FindUserByIDs(ctx, userIds...)
	if err != nil {
		return err
	}
	userMap := make(map[uint64]*User)
	for _, user := range users {
		userMap[user.UID] = user
	}
	for _, status := range statuses {
		var IsPublic bool
		if status.HideTime == nil || status.HideTime.Unix() > time.Now().UTC().Unix() {
			IsPublic = true
		}
		status.IsPublic = IsPublic
		status.User = userMap[status.UID]
	}
	return nil
}

func preloadRelatedStatus(ctx context.Context, statuses ...*Status) error {

	//parent status black
	currentUID, ok := ctx.Value(utils.CurrentUIDKey{}).(uint64)
	blackUids := []uint64{}
	if ok && currentUID > 0 {
		uids, err := AdminListBlackListUserIDs(ctx, currentUID)
		if err == nil {
			blackUids = uids
		}
	}
	statusIds := make([]primitive.ObjectID, 0)
	for _, status := range statuses {
		if !status.ParentID.IsZero() {
			statusIds = append(statusIds, status.ParentID)
		}
		if !status.OriginID.IsZero() {
			statusIds = append(statusIds, status.OriginID)
		}
	}
	relatedStatuses := make([]*Status, 0)
	err := db.ODM(ctx).Where(bson.M{"_id": bson.M{"$in": statusIds}}).Find(&relatedStatuses).Error
	if err != nil {
		return err
	}
	err = PreloadStatusData(ctx, false, relatedStatuses...)
	if err != nil {
		return err
	}
	statusMap := make(map[primitive.ObjectID]*Status)
	for _, status := range relatedStatuses {
		statusMap[status.ID] = status
	}
	for _, status := range statuses {
		if !status.ParentID.IsZero() && statusMap[status.ParentID] == nil {
			status.ParentStatusIsDeleted = true
		}
		status.ParentStatus = statusMap[status.ParentID]
		status.OriginStatus = statusMap[status.OriginID]
		if len(blackUids) > 0 && !status.ParentStatusIsDeleted && status.ParentStatus != nil {
			for _, v := range blackUids {
				if status.ParentStatus.UID == v {
					status.ParentStatusIsBlacked = true
				}
			}
		}
	}
	return nil
}

func preloadAttachment(ctx context.Context, statuses ...*Status) error {
	paths := make([]string, 0)
	linkMetas := make([]*meta.LinkMeta, 0)
	for _, status := range statuses {
		if status.StatusType != enum.LinkStatus {
			continue
		}
		linkMeta := status.LinkMeta
		if linkMeta != nil {
			paths = append(paths, linkMeta.ImagePath)
			linkMetas = append(linkMetas, linkMeta)
		}

	}
	images, err := storage.ImageClient.GetFileUrl(ctx, paths...)
	if err != nil {
		return err
	}
	for _, linkMeta := range linkMetas {
		if images[linkMeta.ImagePath] != "" {
			linkMeta.ImageURL = images[linkMeta.ImagePath]
		}
	}
	return nil
}

func preloadImage(ctx context.Context, statuses ...*Status) error {
	paths := make([]string, 0)
	metas := make([]*meta.ImageMeta, 0)
	for _, status := range statuses {
		if status.StatusType != enum.ImageStatus {
			continue
		}
		meta := status.ImageMeta
		if meta != nil {
			paths = append(paths, meta.Images...)
			metas = append(metas, meta)
		}
	}
	images, err := storage.ImageClient.GetFileUrl(ctx, paths...)
	if err != nil {
		return err
	}
	//thumb image
	opts := &storage.ImageOptions{
		Quality: 15,
	}
	thumbImages, err := storage.ImageClient.GetFileUrlOptions(ctx, opts, paths...)
	for _, meta := range metas {
		meta.ImageURLs = []string{}
		meta.ThumbImageURLs = []string{}
		for _, path := range meta.Images {
			meta.ImageURLs = append(meta.ImageURLs, images[path])
			meta.ThumbImageURLs = append(meta.ThumbImageURLs, thumbImages[path])
		}
	}
	return nil
}
