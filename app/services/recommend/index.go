package recommend

import (
	/*"github.com/mises-id/sns-socialsvc/app/models"
	 "github.com/mises-id/sns-socialsvc/app/services/recommend/data"
	"github.com/mises-id/sns-socialsvc/app/services/recommend/filter"
	"github.com/mises-id/sns-socialsvc/lib/utils" */
	"context"
	"fmt"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/services/recommend/data"
	"github.com/mises-id/sns-socialsvc/app/services/recommend/filter"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ListStatusInput struct {
		UID uint64
		Num uint16
	}
	ListRecommendStatusInput struct {
		UID       uint64
		Num       uint16
		filterIDs []primitive.ObjectID
	}
	ListStarUserStatusInput struct {
		UID       uint64
		Num       uint16
		filterIDs []primitive.ObjectID
	}
)

//list  status ids
func ListStatus(ctx context.Context, in *ListStatusInput) ([]primitive.ObjectID, error) {

	return listStatus(ctx, in)

}

func listStatus(ctx context.Context, in *ListStatusInput) ([]primitive.ObjectID, error) {

	//find recommend status
	statusIDs := make([]primitive.ObjectID, 0)
	recommend_statusIDs, err := listRecommendedStatus(ctx, &ListRecommendStatusInput{
		Num: in.Num,
		UID: in.UID,
	})
	if err != nil {
		return nil, err
	}
	statusIDs = append(statusIDs, recommend_statusIDs...)
	//follow2 status
	return statusIDs, nil
}

//list recommend status ids
func listRecommendedStatus(ctx context.Context, in *ListRecommendStatusInput) ([]primitive.ObjectID, error) {
	if in.Num == 0 {
		return []primitive.ObjectID{}, nil
	}
	num := in.Num
	uid := in.UID
	//get recommend and star user status pool
	statuses, err := data.ListRecommendAndStarUserStatus(ctx)
	if err != nil {
		return nil, err
	}
	statusIDs := make([]primitive.ObjectID, 0)
	for _, status := range statuses {
		if num == 0 {
			break
		}
		statusID := status.ID
		if utils.InArrayObject(statusID, in.filterIDs) > -1 {
			continue
		}
		exist, err := filter.StatusBfExists(ctx, uid, statusID)
		if err != nil {
			fmt.Println("status bf exist error: ", err.Error())
			continue
		}
		if !exist {
			statusIDs = append(statusIDs, statusID)
			num--
		}
	}
	return statusIDs, nil
}

//after list status
func ListStatusAfter(ctx context.Context, uid uint64, statuses []*models.Status) error {

	return listStatusAfter(ctx, uid, statuses)
}

func listStatusAfter(ctx context.Context, uid uint64, statuses []*models.Status) error {
	ids := make([]primitive.ObjectID, len(statuses))
	for k, status := range statuses {
		ids[k] = status.ID
	}
	err := filter.StatusBfInsert(ctx, uid, ids...)
	if err != nil {
		fmt.Println("status bf insert error: ", err.Error())
	}
	return err
}

func InitRecommendData(ctx context.Context) error {
	data.InitStatusGroupUserPool(ctx)
	data.InitStatusRecommend(ctx)
	data.InitStatusStarUserPool(ctx)
	data.InitUserFollowingPool(ctx)
	return nil
}
