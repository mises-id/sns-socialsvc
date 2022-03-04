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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	getListNum              = 20
	userTwitterAuthMaxIdKey = "user_twiter_auth_max_id"
	airdropClient           airdrop.IClient
)

type FaucetCallback struct {
}

func TwitterAirdrop(ctx context.Context) {
	airdropClient.SetListener(&FaucetCallback{})
	airdropCreate(ctx)
	airdropTx(ctx)
	fmt.Println("airdrop finished")
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

func airdropTx(ctx context.Context) {
	airdrops, err := getAirdropList(ctx)
	if err != nil {
		fmt.Println("err: ", err.Error())
		return
	}
	if err := airdropRun(ctx, airdrops); err != nil {
		return
	}
	if len(airdrops) == getListNum {
		airdropTx(ctx)
	}
	return
}

func airdropRun(ctx context.Context, airdrops []*models.Airdrop) error {
	for _, v := range airdrops {
		fmt.Printf("misesid:%s,coin:%d", v.Misesid, v.Coin)
		err := airdropClient.RunAsync(v.Misesid, "", v.Coin)
		if err != nil {
			fmt.Println("airdrop run error: ", err.Error())
			return err
		}
	}
	return nil
}

func (cb *FaucetCallback) OnTxGenerated(cmd types.MisesAppCmd) {
	fmt.Printf("OnTxGenerated\n")
	misesid := cmd.MisesUID()
	txid := cmd.TxID()
	fmt.Println("misesid: ", misesid)
	fmt.Println("tx_id: ", txid)
	err := txGeneratedAfter(context.Background(), misesid, txid)
	if err != nil {
		fmt.Println("tx generated after  error: ", err.Error())
	}

}
func (cb *FaucetCallback) OnSucceed(cmd types.MisesAppCmd) {
	fmt.Printf("OnSucceed\n")
	misesid := cmd.MisesUID()
	txid := cmd.TxID()
	fmt.Println("misesid: ", misesid)
	fmt.Println("tx_id: ", txid)
	err := successAfter(context.Background(), misesid)
	if err != nil {
		fmt.Println("tx success after  error: ", err.Error())
	}

}
func (cb *FaucetCallback) OnFailed(cmd types.MisesAppCmd) {
	fmt.Printf("OnFailed\n")
	misesid := cmd.MisesUID()
	txid := cmd.TxID()
	fmt.Println("misesid: ", misesid)
	fmt.Println("tx_id: ", txid)
	err := failedAfter(context.Background(), misesid)
	if err != nil {
		fmt.Println("tx failed after  error: ", err.Error())
	}

}

func successAfter(ctx context.Context, misesid string) error {
	//airdrop update
	params := &search.AirdropSearch{
		Misesid: misesid,
		Type:    enum.AirdropTwitter,
	}
	airdrop, err := models.FindAirdrop(ctx, params)
	if err != nil {
		fmt.Println("find airdrop error: ", err.Error())
		return err
	}
	if airdrop.Status != 0 {
		fmt.Printf("misesid:%s,  finished", misesid)
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
		Misesid: misesid,
		Type:    enum.AirdropTwitter,
	}
	airdrop, err := models.FindAirdrop(ctx, params)
	if err != nil {
		fmt.Println("find airdrop error: ", err.Error())
		return err
	}
	if airdrop.Status != 0 {
		fmt.Printf("misesid:%s,  finished", misesid)
		return errors.New("misesid finished")
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

func txGeneratedAfter(ctx context.Context, misesid string, tx_id string) error {
	//update
	params := &search.AirdropSearch{
		Misesid: misesid,
		Type:    enum.AirdropTwitter,
	}
	airdrop, err := models.FindAirdrop(ctx, params)
	if err != nil {
		fmt.Println("find airdrop error: ", err.Error())
		return err
	}
	if airdrop.TxID != "" {
		fmt.Printf("misesid:%s has tx_id,old tx_id:%s,new_tx_id:%s", misesid, airdrop.TxID, tx_id)
		return errors.New("tx_id exists")
	}
	//update
	return airdrop.UpdateTxID(ctx, tx_id)
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
	var max, rate float64
	max = 100
	rate = 1000000
	coin := 1 + 0.01*float64(userTwitter.TwitterUser.TweetCount) + 0.005*float64(userTwitter.TwitterUser.FollowersCount)
	if coin > max {
		coin = max
	}
	return int64(coin * rate)
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

func getAirdropList(ctx context.Context) ([]*models.Airdrop, error) {
	params := &search.AirdropSearch{
		NotTxID: true,
		Status:  enum.AirdropDefault,
		ListNum: int64(getListNum),
	}
	return models.ListAirdrop(ctx, params)
}

func SetAirdropClient() {
	airdropClient = airdrop.New()
}
