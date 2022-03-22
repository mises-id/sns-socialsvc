package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/lib/utils"
)

func UserToChain(ctx context.Context) {
	utils.WirteLogDay("./log/user_to_chain.log")
	runUserToChain(ctx)
}

func runUserToChain(ctx context.Context) error {
	lists, err := getChainUser(ctx)
	if err != nil {
		return err
	}
	for _, v := range lists {
		misesClient.SetListener(&RegisterCallback{})
		err1 := misesClient.Register(v.Misesid, v.Pubkey)
		if err1 != nil {
			fmt.Printf("mises[%s] user register chain error:%s \n", v.Misesid, err1.Error())
		}
	}
	return nil
}

func getChainUser(ctx context.Context) ([]*models.ChainUser, error) {

	params := &search.ChainUserSearch{
		Status:   enum.ChainUserDefault,
		SortType: enum.SortAsc,
		SortKey:  "_id",
		ListNum:  20,
	}
	list, err := models.ListChainUser(ctx, params)
	if err != nil {
		fmt.Println("list chain user error: ", err.Error())
		return nil, err
	}
	return list, nil
}

func (cb *RegisterCallback) OnTxGenerated(cmd types.MisesAppCmd) {
	misesid := cmd.MisesUID()
	fmt.Printf("Mises[%s] User Register OnTxGenerated %s\n", misesid, cmd.TxID())
	txGenerated(context.Background(), misesid, cmd.TxID())

}
func (cb *RegisterCallback) OnSucceed(cmd types.MisesAppCmd) {
	misesid := cmd.MisesUID()
	fmt.Printf("Mises[%s] User Register OnSucceed\n", misesid)
	success(context.Background(), misesid)

}
func (cb *RegisterCallback) OnFailed(cmd types.MisesAppCmd, err error) {
	misesid := cmd.MisesUID()
	if err != nil {
		fmt.Printf("Mises[%s] User Register OnFailed: %s\n", misesid, err.Error())
	} else {
		fmt.Printf("Mises[%s] User Register OnFailed\n", misesid)
	}

	failed(context.Background(), misesid)

}

func success(ctx context.Context, misesid string) error {
	//airdrop update
	params := &search.ChainUserSearch{
		Misesid: misesid,
		Status:  enum.ChainUserDefault,
	}

	chainUser, err := models.FindChainUser(ctx, params)
	if err != nil {
		fmt.Println("find chain user error: ", err.Error())
		return err
	}

	if err = chainUser.UpdateStatus(ctx, enum.ChainUserSuccess); err != nil {
		fmt.Println("chain user success update error: ", err.Error())
		return err
	}
	if err = models.UpdateUserOnChainByMisesid(ctx, misesid); err != nil {
		fmt.Println("update user onchain error: ", err.Error())
		chainUser.UpdateStatus(ctx, enum.ChainUserDefault)
		return err
	}
	return nil
}
func failed(ctx context.Context, misesid string) error {

	return nil
}

func txGenerated(ctx context.Context, misesid string, tx_id string) error {
	//update
	params := &search.ChainUserSearch{
		Misesid: misesid,
		Status:  enum.ChainUserDefault,
	}
	chainUser, err := models.FindChainUser(ctx, params)
	if err != nil {
		fmt.Println("find chain user error: ", err.Error())
		return err
	}
	if chainUser.TxID != "" || chainUser.Status != enum.ChainUserDefault {

		return errors.New("chain user tx_id exists")
	}
	//update
	return chainUser.UpdateTxID(ctx, tx_id)
}
