package airdrop

import (
	"context"
	"errors"
	"fmt"

	"time"

	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/lib/airdrop"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	getListNum              = 20
	userTwitterAuthMaxIdKey = "user_twiter_auth_max_id"
	airdropClient           airdrop.IClient
	airdropStop             chan int
)

type FaucetCallback struct {
	ctx context.Context
}

func AirdropTwitter(ctx context.Context) {
	airdropStop = make(chan int)
	utils.WirteLogDay("./log/airdrop.log")
	airdropClient.SetListener(&FaucetCallback{ctx})
	go airdropTx(ctx)
	select {
	case <-airdropStop:
		fmt.Println("airdrop stop")
	}
	return
}

func airdropToStop() {
	airdropStop <- 1
	return
}

func airdropTx(ctx context.Context) {
	airdrop, err := getAirdrop(ctx)
	if err != nil {
		fmt.Println("err: ", err.Error())
		airdropToStop()
		return
	}
	if err := airdropRun(ctx, airdrop); err != nil {
		airdropToStop()
		return
	}
	return
}

func getAirdropList(ctx context.Context) ([]*models.Airdrop, error) {
	params := &search.AirdropSearch{
		NotTxID: true,
		Status:  enum.AirdropDefault,
		ListNum: int64(getListNum),
	}
	return models.ListAirdrop(ctx, params)
}

//get one
func getAirdrop(ctx context.Context) (*models.Airdrop, error) {
	params := &search.AirdropSearch{
		NotTxID:  true,
		SortType: enum.SortAsc,
		SortKey:  "_id",
		Status:   enum.AirdropDefault,
	}
	return models.FindAirdrop(ctx, params)
}

func airdropRun(ctx context.Context, airdrop *models.Airdrop) error {
	fmt.Printf("misesid:%s,coin:%d\n", airdrop.Misesid, airdrop.Coin)
	err := airdropClient.RunAsync(airdrop.Misesid, "", airdrop.Coin)
	if err != nil {
		fmt.Println("airdrop run error: ", err.Error())
		return err
	}
	return pendingAfter(ctx, airdrop.ID)
}

func (cb *FaucetCallback) OnTxGenerated(cmd types.MisesAppCmd) {
	misesid := cmd.MisesUID()
	fmt.Printf("Mises[%s] Airdrop OnTxGenerated %s\n", misesid, cmd.TxID())
	txid := cmd.TxID()
	err := txGeneratedAfter(context.Background(), misesid, txid)
	if err != nil {
		fmt.Println("tx generated after  error: ", err.Error())
	}

}
func (cb *FaucetCallback) OnSucceed(cmd types.MisesAppCmd) {
	misesid := cmd.MisesUID()
	fmt.Printf("Mises[%s] Airdrop OnSucceed\n", misesid)
	err := successAfter(context.Background(), misesid)
	if err != nil {
		fmt.Println("tx success after  error: ", err.Error())
	}
	airdropTx(cb.ctx)
}

func (cb *FaucetCallback) OnFailed(cmd types.MisesAppCmd, err error) {
	misesid := cmd.MisesUID()
	if err != nil {
		fmt.Printf("Mises[%s] Airdrop OnFailed: %s\n", misesid, err.Error())
	}
	err = failedAfter(context.Background(), misesid)
	if err != nil {
		fmt.Println("tx failed after  error: ", err.Error())
	}

}

func successAfter(ctx context.Context, misesid string) error {
	//airdrop update
	params := &search.AirdropSearch{
		Misesid: misesid,
		Type:    enum.AirdropTwitter,
		Status:  enum.AirdropPending,
	}
	airdrop, err := models.FindAirdrop(ctx, params)
	if err != nil {
		fmt.Println("find airdrop error: ", err.Error())
		return err
	}
	if airdrop.Status != enum.AirdropPending {
		fmt.Printf("misesid:%s,  status error ", misesid)
		return errors.New("misesid finished")
	}
	if err = airdrop.UpdateStatus(ctx, enum.AirdropSuccess); err != nil {
		fmt.Println("airdrop success update error: ", err.Error())
		return err
	}
	//update user airdrop coin
	if err = updateUserAirdrop(ctx, airdrop.UID, uint64(airdrop.Coin)); err != nil {
		fmt.Println("airdrop success update user ext error: ", err.Error())
		return err
	}
	return nil
}
func failedAfter(ctx context.Context, misesid string) error {
	//airdrop update
	params := &search.AirdropSearch{
		Misesid:  misesid,
		Type:     enum.AirdropTwitter,
		Statuses: []enum.AirdropStatus{enum.AirdropDefault, enum.AirdropPending},
	}
	airdrop, err := models.FindAirdrop(ctx, params)
	if err != nil {
		fmt.Println("find airdrop error: ", err.Error())
		return err
	}
	if airdrop.Status != enum.AirdropPending && airdrop.Status != enum.AirdropDefault {
		fmt.Printf("misesid:%s,  status error", misesid)
		return errors.New("airdrop status error")
	}
	if err = airdrop.UpdateStatus(ctx, enum.AirdropFailed); err != nil {
		fmt.Println("airdrop failed update error: ", err.Error())
		return err
	}
	return nil
}

func updateUserAirdrop(ctx context.Context, uid uint64, coin uint64) error {
	user_ext, err := models.FindOrCreateUserExt(ctx, uid)
	if err != nil {
		return err
	}
	user_ext.AirdropCoin += coin

	return user_ext.UpdateAirdrop(ctx)
}

func pendingAfter(ctx context.Context, id primitive.ObjectID) error {
	params := &search.AirdropSearch{
		ID:     id,
		Type:   enum.AirdropTwitter,
		Status: enum.AirdropDefault,
	}
	airdrop, err := models.FindAirdrop(ctx, params)
	if err != nil {
		fmt.Println("pending after find airdrop error: ", err.Error())
		return err
	}
	if airdrop.TxID != "" || airdrop.Status != enum.AirdropDefault {
		return errors.New("pending status tx_id exists")
	}
	return airdrop.UpdateStatusPending(ctx)
}

func txGeneratedAfter(ctx context.Context, misesid string, tx_id string) error {
	//update
	params := &search.AirdropSearch{
		Misesid: misesid,
		Type:    enum.AirdropTwitter,
		Status:  enum.AirdropPending,
	}
	airdrop, err := models.FindAirdrop(ctx, params)
	if err != nil {
		fmt.Println("find airdrop error: ", err.Error())
		return err
	}
	if airdrop.TxID != "" || airdrop.Status != enum.AirdropPending {
		return errors.New("tx_id exists")
	}
	//update
	return airdrop.UpdateTxID(ctx, tx_id)
}

func CretaeAirdropTwitter(ctx context.Context) {

	utils.WirteLogDay("./log/create_airdrop.log")
	if !models.GetAirdropStatus(ctx) {
		return
	}
	airdropCreate(ctx)
}

func airdropCreate(ctx context.Context) {
	list, err := createdTwitterAirdrop(ctx)
	if err != nil {
		return
	}
	if len(list) == getListNum {
		airdropCreate(ctx)
	}
}

func createdTwitterAirdrop(ctx context.Context) ([]*models.Airdrop, error) {
	//get user twitter auth
	userTwitterAuthList, err := getAirdropUserTwitterAuth(ctx)
	if err != nil {
		return nil, err
	}
	num := len(userTwitterAuthList)
	if num == 0 {
		return []*models.Airdrop{}, nil
	}
	airdrops := make([]*models.Airdrop, 0)
	for _, v := range userTwitterAuthList {
		airdrop := &models.Airdrop{
			UID:       v.UID,
			Misesid:   v.Misesid,
			Status:    enum.AirdropDefault,
			Type:      enum.AirdropTwitter,
			Coin:      getTwitterAirdropCoin(ctx, v),
			TxID:      "",
			CreatedAt: time.Now(),
		}
		airdrops = append(airdrops, airdrop)
	}
	err1 := models.CreateAirdropMany(ctx, airdrops)
	if err1 != nil {
		fmt.Println("create airdrop error: ", err1.Error())
		return nil, err
	}
	maxId := userTwitterAuthList[num-1].ID
	//update
	if err := updateMaxId(ctx, maxId); err != nil {
		fmt.Println("update maxid error: ", err.Error())
		return nil, err
	}
	return airdrops, nil
}

func getTwitterAirdropCoin(ctx context.Context, userTwitter *models.UserTwitterAuth) int64 {

	var max, umises, mises uint64
	umises = 1
	mises = 1000000 * umises
	max = 100 * mises
	coin := mises + 10000*umises*userTwitter.TwitterUser.TweetCount + 500*umises*userTwitter.TwitterUser.FollowersCount
	if coin > max {
		coin = max
	}
	return int64(coin)
}

func getMaxId(ctx context.Context) primitive.ObjectID {
	c, err := models.FindOrCreateConfig(ctx, userTwitterAuthMaxIdKey, primitive.NilObjectID)
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
	return models.UpdateOrCreateConfig(ctx, userTwitterAuthMaxIdKey, value)
}

func getAirdropUserTwitterAuth(ctx context.Context) ([]*models.UserTwitterAuth, error) {

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

func SetAirdropClient() {
	airdropClient = airdrop.New()
}
