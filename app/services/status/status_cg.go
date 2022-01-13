package status

import (
	"context"
	"fmt"
	"time"

	"github.com/mises-id/sns-socialsvc/admin"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
)

var (
	updateUserCursor   *models.UserExt
	newRecommendInput  *NewRecommendInput
	newRecommendOutput *NewRecommendOutput
)

type (
	NewRecommendInput struct {
		LastRecommendTime int64
		LastCommonTime    int64
	}
	NewRecommendNext struct {
		LastRecommendTime int64
		LastCommonTime    int64
	}
	NewRecommendOutput struct {
		Data []*models.Status
		Next *NewRecommendNext
	}
)

// new recommend status
func NewRecommendStatus(ctx context.Context, uid uint64, in *NewRecommendInput) (*NewRecommendOutput, error) {

	var totalNum, following2Num, recommendPoolNum, commonPoolNum int64
	//start
	updateUserCursor = &models.UserExt{
		UID: uid,
	}
	newRecommendInput = in
	newRecommendOutput = &NewRecommendOutput{
		Next: &NewRecommendNext{LastRecommendTime: 0, LastCommonTime: 0},
	}
	totalNum = 10
	//following2 pool status
	following2Num = 5
	following2_status_list, err := findListFollowing2Status(ctx, uid, following2Num)
	if err != nil {
		return nil, err
	}
	//recommend pool status
	now_following2_num := len(following2_status_list)
	recommendPoolNum = 5 + following2Num - int64(now_following2_num)
	//TODO filter problem user status && filter black list
	recommend_pool_status_list, err := findListRecommendStatus(ctx, uid, recommendPoolNum)
	if err != nil {
		return nil, err
	}
	//common pool status
	now_recommend_num := len(recommend_pool_status_list)
	commonPoolNum = totalNum - int64(now_following2_num+now_recommend_num)
	common_pool_status, err := findListCommonStatus(ctx, uid, commonPoolNum)
	if err != nil {
		return nil, err
	}
	now_common_num := len(common_pool_status)
	//now_total_num := now_following2_num + now_recommend_num + now_comment_num
	fmt.Printf("following2_num:%d,recommend_num:%d,common_num:%d", now_following2_num, now_recommend_num, now_common_num)
	data := append(following2_status_list, append(recommend_pool_status_list, common_pool_status...)...)
	newRecommendOutput.Data = data
	if newRecommendOutput.Next.LastRecommendTime == 0 {
		newRecommendOutput.Next.LastRecommendTime = in.LastRecommendTime
	}
	if newRecommendOutput.Next.LastCommonTime == 0 {
		newRecommendOutput.Next.LastCommonTime = in.LastCommonTime
	}
	//end update cursor
	if uid > 0 {
		updateUserCursor.Update(ctx)
	}
	return newRecommendOutput, err
}

//find following2 pool status
func findListFollowing2Status(ctx context.Context, uid uint64, num int64) ([]*models.Status, error) {

	if uid == 0 || num <= 0 {
		return []*models.Status{}, nil
	}
	//find following2 uids
	uids, err := findUserFollowing2Uids(ctx, uid)
	if err != nil {
		return nil, err
	}
	// follow2 empty
	if len(uids) == 0 {
		return []*models.Status{}, nil
	}
	start_time := time.Now().AddDate(0, 0, -7)
	params := &admin.AdminStatusParams{
		NInTags:   []enum.TagType{enum.TagRecommendStatus},
		ListNum:   num,
		UIDs:      uids,
		OnlyShow:  true,
		StartTime: &start_time,
		SortType:  1,
		SortKey:   "created_at",
		FromTypes: []enum.FromType{enum.FromForward, enum.FromPost, enum.FromComment},
	}
	//find status recommend pool cursor
	cursors := getUserFollowing2Cursor(ctx, uid)
	if cursors != nil {
		params.ScoreMax = cursors.Max
		params.ScoreMin = cursors.Min
	}

	status_list, err := models.AdminListStatus(ctx, params)
	if err != nil {
		return nil, err
	}
	status_num := len(status_list)
	//update pool status cursor
	if status_num > 0 {
		max, min := getStatusListScoreMaxMin(status_list)
		updateUserFollowing2Cursor(ctx, uid, cursors, max, min)
	}
	return status_list, nil
}

//find recommend pool status
func findListRecommendStatus(ctx context.Context, uid uint64, num int64) ([]*models.Status, error) {

	if num <= 0 {
		return []*models.Status{}, nil
	}
	var err error
	var recommend_pool_cursors *models.RecommendStatusPoolCursor
	recommend_start_time := time.Now().AddDate(0, 0, -14)
	recommend_pool_status_params := &admin.AdminStatusParams{
		Tag:       enum.TagRecommendStatus,
		ListNum:   num,
		OnlyShow:  true,
		StartTime: &recommend_start_time,
		SortType:  1,
		SortKey:   "created_at",
		FromTypes: []enum.FromType{enum.FromForward, enum.FromPost},
	}
	if uid > 0 {
		recommend_pool_cursors = getUserRecommendCursor(ctx, uid)
		if recommend_pool_cursors != nil {
			recommend_pool_status_params.ScoreMax = recommend_pool_cursors.Max
			recommend_pool_status_params.ScoreMin = recommend_pool_cursors.Min
		}
	} else {
		//not login
		recommend_pool_status_params.SortType = -1
		smax := time.Now().UnixMilli()
		if newRecommendInput != nil && newRecommendInput.LastRecommendTime > 0 {
			smax = newRecommendInput.LastRecommendTime
		}
		recommend_pool_status_params.ScoreMax = smax
	}
	recommend_pool_status_list, err := models.AdminListStatus(ctx, recommend_pool_status_params)
	if err != nil {
		return nil, err
	}
	//recommend_status_num := len(recommend_pool_status_list)
	max, min := getStatusListScoreMaxMin(recommend_pool_status_list)
	//update recommend pool status cursor
	if uid > 0 {
		updateUserRecommendCursor(ctx, uid, recommend_pool_cursors, max, min)
	} else {
		newRecommendOutput.Next.LastRecommendTime = min
	}

	return recommend_pool_status_list, nil
}

//find common pool status
func findListCommonStatus(ctx context.Context, uid uint64, num int64) ([]*models.Status, error) {

	if num <= 0 {
		return []*models.Status{}, nil
	}
	var cursors *models.CommonPoolCursor
	start_time := time.Now().AddDate(0, 0, -7)
	params := &admin.AdminStatusParams{
		NInTags:   []enum.TagType{enum.TagRecommendStatus}, //filter recommend status
		ListNum:   num,
		OnlyShow:  true,
		StartTime: &start_time,
		SortType:  1,
		SortKey:   "created_at",
		FromTypes: []enum.FromType{enum.FromForward, enum.FromPost},
	}
	//login user
	if uid > 0 {
		uids, err := findUserFollowing2Uids(ctx, uid)

		if err == nil && len(uids) > 0 {
			fmt.Println("uids:", uids)
			params.NInUIDs = uids //filter following2 user status
		}
		//find status recommend pool cursor
		cursors = getUserCommonCursor(ctx, uid)
		if cursors != nil {
			params.ScoreMax = cursors.Max
			params.ScoreMin = cursors.Min
		}
	} else {
		//not login
		params.SortType = -1
		smax := time.Now().UnixMilli()
		if newRecommendInput != nil && newRecommendInput.LastCommonTime > 0 {
			smax = newRecommendInput.LastCommonTime
		}
		params.ScoreMax = smax
	}

	status_list, err := models.AdminListStatus(ctx, params)
	if err != nil {
		return nil, err
	}
	max, min := getStatusListScoreMaxMin(status_list)
	//update recommend pool status cursor
	if uid > 0 {
		updateUserCommonCursor(ctx, uid, cursors, max, min)
	} else {
		newRecommendOutput.Next.LastCommonTime = min
	}
	return status_list, nil
}

//get status list min max
func getStatusListScoreMaxMin(statuses []*models.Status) (max int64, min int64) {
	status_num := len(statuses)
	if status_num == 0 {
		return max, min
	}
	for _, status := range statuses {
		cmt := status.CreatedAt.UnixMilli()
		if cmt > max {
			max = cmt
		}
		if min == 0 || cmt < min {
			min = cmt
		}
	}
	return max, min
}

//find user following following uids
func findUserFollowing2Uids(ctx context.Context, uid uint64) ([]uint64, error) {

	followingUids, err := models.AdminListFollowingUserIDs(ctx, []uint64{uid})
	if err != nil {
		fmt.Println("find user followingUids error: ", err.Error())
		return nil, err
	}
	uids, err := models.AdminListFollowingUserIDs(ctx, followingUids)

	if err != nil {
		fmt.Println("find user following following uids error: ", err.Error())
	}

	return uids, nil
}

//get user following2 cursor
func getUserFollowing2Cursor(ctx context.Context, uid uint64) *models.Following2PoolCursor {

	//cursor := &models.Following2PoolCursor{Max: 0, Min: 0}
	user_ext, err := models.FindOrCreateUserExt(ctx, uid)
	if err != nil {
		fmt.Println("find or create user ext error: ", err.Error())
		return nil
	}
	return user_ext.Following2PoolCursor

}

//get user recommend cursor
func getUserRecommendCursor(ctx context.Context, uid uint64) *models.RecommendStatusPoolCursor {

	//cursor := &models.RecommendStatusPoolCursor{Max: 0, Min: 0}
	user_ext, err := models.FindOrCreateUserExt(ctx, uid)
	if err != nil {
		fmt.Println("find or create user ext error: ", err.Error())
		return nil
	}
	return user_ext.RecommendStatusPoolCursor

}

//get user common cursor
func getUserCommonCursor(ctx context.Context, uid uint64) *models.CommonPoolCursor {

	//cursor := &models.CommonPoolCursor{Max: 0, Min: 0}
	user_ext, err := models.FindOrCreateUserExt(ctx, uid)
	if err != nil {
		fmt.Println("find or create user ext error: ", err.Error())
		return nil
	}
	return user_ext.CommonPoolCursor

}

//update user following2 cursor
func updateUserFollowing2Cursor(ctx context.Context, uid uint64, pool_cursors *models.Following2PoolCursor, max, min int64) {
	//init
	if pool_cursors == nil || pool_cursors.Max == 0 || pool_cursors.Min == 0 {
		pool_cursors = &models.Following2PoolCursor{}
		pool_cursors.Min = min
		pool_cursors.Max = max
	}
	//update min
	if pool_cursors.Min > min {
		pool_cursors.Min = min
	}
	//update max
	if pool_cursors.Max < max {
		pool_cursors.Max = max
	}

	updateUserCursor.Following2PoolCursor = pool_cursors
}

//update user recommend cursor
func updateUserRecommendCursor(ctx context.Context, uid uint64, pool_cursors *models.RecommendStatusPoolCursor, max, min int64) {
	//init
	if pool_cursors == nil || pool_cursors.Max == 0 || pool_cursors.Min == 0 {
		pool_cursors = &models.RecommendStatusPoolCursor{}
		pool_cursors.Min = min
		pool_cursors.Max = max
	}
	//update min
	if pool_cursors.Min > min {
		pool_cursors.Min = min
	}
	//update max
	if pool_cursors.Max < max {
		pool_cursors.Max = max
	}

	updateUserCursor.RecommendStatusPoolCursor = pool_cursors
}

//update user common cursor
func updateUserCommonCursor(ctx context.Context, uid uint64, pool_cursors *models.CommonPoolCursor, max, min int64) {
	//init
	if pool_cursors == nil || pool_cursors.Max == 0 || pool_cursors.Min == 0 {
		pool_cursors = &models.CommonPoolCursor{}
		pool_cursors.Min = min
		pool_cursors.Max = max
	}
	//update min
	if pool_cursors.Min > min {
		pool_cursors.Min = min
	}
	//update max
	if pool_cursors.Max < max {
		pool_cursors.Max = max
	}

	updateUserCursor.CommonPoolCursor = pool_cursors
}
