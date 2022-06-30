package channel_user

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/app/services/user_twitter"
	airdropLib "github.com/mises-id/sns-socialsvc/lib/airdrop"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	getListNum                 = 10
	channelTwitterAuthMaxIdKey = "channel_twiter_auth_max_id"
	airdropStop                chan int
	airdropDo                  bool
	totalAirdropNum            int
)

type (
	FaucetCallback struct {
		ctx context.Context
	}
)

//airdrop channel
func AirdropChannel(ctx context.Context) {
	totalAirdropNum = 100
	airdropStop = make(chan int)
	airdropDo = true
	fmt.Println("airdrop channel start")
	airdropLib.AirdropClient.SetListener(&FaucetCallback{ctx})
	go airdropTx(ctx)
	select {
	case <-airdropStop:
		fmt.Println("airdrop channel stop")
	}
	return
}

func airdropToStop() {
	airdropDo = false
	airdropStop <- 1
	return
}

func airdropTx(ctx context.Context) {
	airdrops, err := getChannelAirdropList(ctx)
	if err != nil {
		airdropToStop()
		return
	}
	for _, airdrop := range airdrops {
		if err := airdropRun(ctx, airdrop); err != nil {
			airdropToStop()
			return
		}
	}
	return
}

func airdropTxOne(ctx context.Context) {
	airdrop, err := getChannelAirdrop(ctx)
	if err != nil {
		airdropToStop()
		return
	}
	if err := airdropRun(ctx, airdrop); err != nil {
		airdropToStop()
		return
	}
	return
}

func getChannelAirdropList(ctx context.Context) ([]*models.ChannelUser, error) {
	params := &search.ChannelUserSearch{
		ValidStates:   []enum.UserValidState{enum.UserValidSucessed},
		SortType:      enum.SortAsc,
		SortKey:       "_id",
		AirdropStates: []enum.ChannelAirdropState{enum.ChannelAirdropDefault},
		ListNum:       int64(getListNum),
	}
	return models.ListChannelUser(ctx, params)
}

//get one
func getChannelAirdrop(ctx context.Context) (*models.ChannelUser, error) {
	params := &search.ChannelUserSearch{
		ValidStates:   []enum.UserValidState{enum.UserValidSucessed},
		SortType:      enum.SortAsc,
		SortKey:       "_id",
		AirdropStates: []enum.ChannelAirdropState{enum.ChannelAirdropDefault},
	}
	return models.FindChannelUser(ctx, params)
}

func airdropRun(ctx context.Context, channel_user *models.ChannelUser) error {
	if totalAirdropNum <= 0 {
		return errors.New("too many airdrop num")
	}
	misesid := channel_user.ChannelMisesid
	amount := channel_user.Amount
	trackid := channel_user.ID.Hex()
	fmt.Printf("channel airdrop num:%d,id:%s,coin:%d\n", totalAirdropNum, trackid, amount)
	err := airdropLib.AirdropClient.RunAsync(misesid, "", amount, airdropLib.AirdropClient.SetTrackID(trackid))
	if err != nil {
		return err
	}
	totalAirdropNum--
	return pendingAfter(ctx, channel_user.ID)
}

func trackIDToObjectID(trackid string) primitive.ObjectID {

	id, err := primitive.ObjectIDFromHex(trackid)
	if err != nil {
		fmt.Println("trackid error: ", err.Error())
		id = primitive.NilObjectID
	}
	return id
}

func (cb *FaucetCallback) OnTxGenerated(cmd types.MisesAppCmd) {
	trackid := cmd.TrackID()
	id := trackIDToObjectID(trackid)
	fmt.Printf("ID[%s] Channel Airdrop OnTxGenerated %s\n", trackid, cmd.TxID())
	txid := cmd.TxID()
	err := txGeneratedAfter(context.Background(), id, txid)
	if err != nil {
		fmt.Printf("ID[%s], channel airdrop tx generated after  error:%s \n ", trackid, err.Error())
	}
}
func (cb *FaucetCallback) OnSucceed(cmd types.MisesAppCmd) {
	txid := cmd.TxID()
	trackid := cmd.TrackID()
	id := trackIDToObjectID(trackid)
	fmt.Printf("ID[%s] Channel Airdrop OnSucceed %s\n", trackid, cmd.TxID())
	err := successAfter(context.Background(), id)
	if err != nil {
		fmt.Printf("ID[%s],TxID[%s] ,channel airdrop tx success after  error:%s \n", trackid, txid, err.Error())
	}
	if airdropDo {
		airdropTxOne(cb.ctx)
	}
}

func (cb *FaucetCallback) OnFailed(cmd types.MisesAppCmd, err error) {
	txid := cmd.TxID()
	trackid := cmd.TrackID()
	id := trackIDToObjectID(trackid)
	if err != nil {
		fmt.Printf("ID[%s],TxID[%s], Channel Airdrop OnFailed: %s\n", trackid, txid, err.Error())
	}
	err = failedAfter(context.Background(), id, err.Error())
	if err != nil {
		fmt.Printf("ID[%s],TxID[%s] ,channel airdrop tx failed after  error:%s \n", trackid, txid, err.Error())
	}
	if airdropDo {
		airdropTxOne(cb.ctx)
	}
}

func successAfter(ctx context.Context, id primitive.ObjectID) error {
	//airdrop update
	/* params := &search.ChannelUserSearch{
		ID: id,
	} */
	channel_user, err := models.FindChannelUserByID(ctx, id)
	if err != nil {
		fmt.Println("channel airdrop find channel user error: ", err.Error())
		return err
	}
	if channel_user.AirdropState != enum.ChannelAirdropPending {
		fmt.Printf("id:%s,  state not pending, error ", id)
		return errors.New("state error")
	}
	if err = channel_user.UpdateStatusSuccess(ctx); err != nil {
		fmt.Println("channel airdrop success update error: ", err.Error())
		return err
	}
	//update user airdrop coin
	if err = updateUserAirdrop(ctx, channel_user.ChannelUID, channel_user.Amount); err != nil {
		fmt.Println("channel airdrop success update user ext error: ", err.Error())
		return err
	}
	return nil
}
func failedAfter(ctx context.Context, id primitive.ObjectID, airdrop_err string) error {
	//airdrop update
	/* params := &search.ChannelUserSearch{
		ID: id,
	} */
	channel_user, err := models.FindChannelUserByID(ctx, id)
	if err != nil {
		fmt.Println("channel airdrop find channel user error: ", err.Error())
		return err
	}
	if channel_user.AirdropState != enum.ChannelAirdropPending {
		fmt.Printf("id:%s,  state error", id)
		return errors.New("channel airdrop state error")
	}
	if err = channel_user.UpdateStatusFailed(ctx, airdrop_err); err != nil {
		fmt.Println("channel airdrop failed update error: ", err.Error())
		return err
	}
	return nil
}

func updateUserAirdrop(ctx context.Context, uid uint64, coin int64) error {
	user_ext, err := models.FindOrCreateUserExt(ctx, uid)
	if err != nil {
		return err
	}
	user_ext.ChannelAirdropCoin += uint64(coin)

	return user_ext.UpdateChannelAirdrop(ctx)
}

func pendingAfter(ctx context.Context, id primitive.ObjectID) error {
	/* params := &search.ChannelUserSearch{
		ID: id,
	} */
	channel_user, err := models.FindChannelUserByID(ctx, id)
	if err != nil {
		fmt.Printf("id[%s],channel airdrop pending after find channel user error: %s \n", id.Hex(), err.Error())
		return err
	}
	if channel_user.TxID != "" && channel_user.AirdropState != enum.ChannelAirdropDefault {
		return errors.New("channel airdrop pending state tx_id exists")
	}
	err = channel_user.UpdateStatusPending(ctx)
	if err != nil {
		fmt.Println("cg: ", err.Error())
	}
	return err
}

func txGeneratedAfter(ctx context.Context, id primitive.ObjectID, tx_id string) error {
	//update
	/* params := &search.ChannelUserSearch{
		ID: id,
	} */
	channel_user, err := models.FindChannelUserByID(ctx, id)
	if err != nil {
		fmt.Println("channel airdrop find  error: ", err.Error())
		return err
	}
	if channel_user.TxID != "" || channel_user.AirdropState != enum.ChannelAirdropPending {
		return errors.New("tx_id exists")
	}
	//update
	return channel_user.UpdateTxID(ctx, tx_id)
}

//create channel airdrop
func CretaeChannelAirdrop(ctx context.Context) error {

	utils.WirteLogDay("./log/create_channel_airdrop.log")
	if !models.GetAirdropStatus(ctx) {
		return nil
	}
	return channelAirdropCreate(ctx)
}

func channelAirdropCreate(ctx context.Context) error {
	c, err := countUserTwitterAuth(ctx)
	if err != nil {
		return err
	}
	if c == 0 {
		return nil
	}
	times := int(math.Ceil(float64(c) / float64(getListNum)))
	for i := 0; i < times; i++ {
		err := createdChannelAirdrop(ctx)
		if err != nil {
			fmt.Println("create channel airdrop error: ", err.Error())
			return err
		}
	}
	return nil
}

func createdChannelAirdrop(ctx context.Context) error {
	//get user twitter auth
	userTwitterAuthList, err := getUserTwitterAuth(ctx)
	if err != nil {
		return err
	}
	num := len(userTwitterAuthList)
	if num == 0 {
		return nil
	}
	uids := make([]uint64, 0)
	userTwitterAuthMap := make(map[uint64]*models.UserTwitterAuth, num)
	for _, v := range userTwitterAuthList {
		uids = append(uids, v.UID)
		userTwitterAuthMap[v.UID] = v
	}
	channel_users, err := getNotAuthChannelUserByUIDs(ctx, uids...)
	if err != nil {
		return err
	}
	if len(channel_users) > 0 {
		for _, channel_user := range channel_users {
			var amount int64
			var valid_state enum.UserValidState
			valid_state = enum.UserValidFailed
			//valid twitter register time
			twitter_user := userTwitterAuthMap[channel_user.UID]
			is_valid := user_twitter.IsValidTwitterUser(twitter_user.TwitterUser)
			if is_valid {
				amount = getChannelAirdropCoin(ctx, twitter_user)
				valid_state = enum.UserValidSucessed
			}
			err := channel_user.UpdateCreateAirdrop(ctx, valid_state, amount)
			if err != nil {
				fmt.Printf("uid[%d], update create airdrop error:%s \n", channel_user.UID, err.Error())
				return err
			}
		}
	}
	maxId := userTwitterAuthList[num-1].ID
	//update
	if err := updateMaxId(ctx, maxId); err != nil {
		fmt.Println("update maxid error: ", err.Error())
		return err
	}
	return nil
}

func getNotAuthChannelUserByUIDs(ctx context.Context, uids ...uint64) ([]*models.ChannelUser, error) {
	params := &search.ChannelUserSearch{
		UIDs:        uids,
		ValidStates: []enum.UserValidState{enum.UserValidDefalut, enum.UserValidFailed},
	}
	return models.ListChannelUser(ctx, params)
}

func getChannelAirdropCoin(ctx context.Context, userTwitter *models.UserTwitterAuth) int64 {
	if userTwitter.TwitterUser.FollowersCount == 0 {
		return 0
	}
	var max, umises, mises uint64
	umises = 1
	mises = 1000000 * umises
	max = 100 * mises
	/* tweet_count := userTwitter.TwitterUser.TweetCount
	if tweet_count > 500 {
		tweet_count = 500
	} */
	//coin := mises + 10000*umises*tweet_count + 5000*umises*userTwitter.TwitterUser.FollowersCount
	coin := mises + 5000*umises*userTwitter.TwitterUser.FollowersCount
	if coin > max {
		coin = max
	}
	return int64(coin) / 10
}

func getMaxId(ctx context.Context) primitive.ObjectID {
	c, err := models.FindOrCreateConfig(ctx, channelTwitterAuthMaxIdKey, primitive.NilObjectID)
	if err != nil {
		fmt.Println("find or create config error: ", err.Error())
		return primitive.NilObjectID
	}
	gid := c.Value
	return gid.(primitive.ObjectID)
}

func updateMaxId(ctx context.Context, max_id primitive.ObjectID) error {
	var value interface{}
	value = max_id
	return models.UpdateOrCreateConfig(ctx, channelTwitterAuthMaxIdKey, value)
}

func getUserTwitterAuth(ctx context.Context) ([]*models.UserTwitterAuth, error) {

	params := &search.UserTwitterAuthSearch{
		GID:      getMaxId(ctx),
		SortType: enum.SortAsc,
		SortKey:  "_id",
		ListNum:  int64(getListNum),
	}
	list, err := models.ListUserTwitterAuth(ctx, params)
	if err != nil {
		fmt.Println("list user twitter auth error: ", err.Error())
		return nil, err
	}
	return list, nil
}
func countUserTwitterAuth(ctx context.Context) (int64, error) {

	params := &search.UserTwitterAuthSearch{
		GID:      getMaxId(ctx),
		SortType: enum.SortAsc,
		SortKey:  "_id",
		ListNum:  int64(getListNum),
	}
	c, err := models.CountUserTwitterAuth(ctx, params)
	if err != nil {
		fmt.Println("count user twitter auth error: ", err.Error())
		return c, err
	}
	return c, nil
}
