package nft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/config/env"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/sirupsen/logrus"
)

type (
	SingleAssetInput struct {
		AssetContractAddress string
		TokenId              string
		AccountAddress       string
		IncludeOrders        string
		Network              string
	}
	AssetContractInput struct {
		AssetContractAddress string
		Network              string
	}
	OpensaeInput struct {
		AssetContractAddress string
		TokenId              string
		Owner                string
		Limit                uint64
		Cursor               string
		Network              string
	}
	ListAssetInput struct {
		Owner   string
		Limit   uint64
		Cursor  string
		Network string
	}
	ListAssetOutput struct {
		Assets   []*models.Asset `json:"assets"`
		Next     string          `json:"next"`
		Previous string          `json:"previous"`
	}
	ListEventOutput struct {
		AssetEvents []*models.AssetEvent `json:"asset_events"`
		Next        string               `json:"next"`
		Previous    string               `json:"previous"`
	}
	HttpResult struct {
		body       []byte
		status     string // e.g. "200 OK"
		statusCode int    // e.g. 200
	}
)

func (res *HttpResult) Restult(out interface{}) error {
	if res.statusCode != http.StatusOK {
		return errors.New(res.status)
	}
	json.Unmarshal(res.body, out)
	return nil
}

func (res *HttpResult) String() string {

	return string(res.body)
}

var (
	openseaApiUrl     = "https://api.opensea.io/api/v1/"
	openseaTestApiUrl = "https://testnets-api.opensea.io/api/v1/"
	defaultTokenID    string
	xApiKey           string
)

func init() {
	defaultTokenID = "1"
	xApiKey = env.Envs.OpenseaApiKey
}

func ListAsset(ctx context.Context, currentUID uint64, in *ListAssetInput) (string, error) {
	user, err := models.FindUserEthAddress(ctx, currentUID)
	if err != nil {
		return "", err
	}
	in.Owner = user.EthAddress
	out, err := ListAssetApi(ctx, in)
	if err != nil {
		return "", codes.ErrTooManyRequests.Newf(err.Error())
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		db.SetupMongo(ctx)
		err2 := AfterListAsset(ctx, currentUID, out)
		if err != nil {
			logrus.Errorln("after_list_asset err: ", err2.Error())
		}
	}()
	return out.String(), nil
}
func AfterListAsset(ctx context.Context, currentUID uint64, httpResult *HttpResult) error {
	out := &ListAssetOutput{}
	httpResult.Restult(out)
	err := SaveUserNftAsset(ctx, currentUID, out.Assets)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}
func ListAssetOut(ctx context.Context, in *ListAssetInput) (*ListAssetOutput, error) {
	res := &ListAssetOutput{}
	out, err := ListAssetApi(ctx, in)
	if err != nil {
		return nil, codes.ErrTooManyRequests.Newf(err.Error())
	}
	out.Restult(res)
	return res, nil
}

func ListAssetApi(ctx context.Context, in *ListAssetInput) (*HttpResult, error) {
	if in.Owner == "" {
		return nil, codes.ErrInvalidArgument.Newf("invalid owner params")
	}
	if in.Limit >= 50 || in.Limit <= 0 {
		in.Limit = 50
	}
	api := openseaApiUrl
	cursorKey := "cursor"
	if in.Network == "test" {
		api = openseaTestApiUrl
		//cursorKey = "offset"
	}
	queryParams := "assets/?owner=" + in.Owner + "&limit=" + strconv.Itoa(int(in.Limit))
	if in.Cursor != "" {
		queryParams = queryParams + "&" + cursorKey + "=" + in.Cursor
	}

	apiUrl := api + queryParams
	return doOpenseaApi(ctx, apiUrl, in.Network)

}
func ListEventOut(ctx context.Context, in *OpensaeInput) (*ListEventOutput, error) {
	res := &ListEventOutput{}
	out, err := ListEventApi(ctx, in)
	if err != nil {
		return nil, codes.ErrTooManyRequests.Newf(err.Error())
	}
	out.Restult(res)
	return res, nil
}
func ListEventApi(ctx context.Context, in *OpensaeInput) (*HttpResult, error) {
	if in.Limit >= 50 || in.Limit <= 0 {
		in.Limit = 50
	}
	api := openseaApiUrl
	if in.Network == "test" {
		api = openseaTestApiUrl
	}
	queryParams := "events/?limit=" + strconv.Itoa(int(in.Limit))
	if in.Cursor != "" {
		queryParams = queryParams + "&cursor" + "=" + in.Cursor
	}
	if in.AssetContractAddress != "" {
		queryParams = queryParams + "&asset_contract_address" + "=" + in.AssetContractAddress
	}
	if in.TokenId != "" {
		queryParams = queryParams + "&token_id" + "=" + in.TokenId
	}

	apiUrl := api + queryParams
	return doOpenseaApi(ctx, apiUrl, in.Network)

}

func GetAssetContract(ctx context.Context, in *AssetContractInput) (string, error) {
	if in.AssetContractAddress == "" {
		return "", codes.ErrInvalidArgument.Newf("invalid asset_contract_address params")
	}

	queryParams := "asset_contract/" + in.AssetContractAddress
	api := openseaApiUrl
	if in.Network == "test" {
		api = openseaTestApiUrl
	}
	apiUrl := api + queryParams
	out, err := doOpenseaApi(ctx, apiUrl, in.Network)
	if err != nil {
		return "", codes.ErrTooManyRequests.Newf(err.Error())
	}
	return out.String(), nil
}

func GetSingleAssetApi(ctx context.Context, in *SingleAssetInput) (*HttpResult, error) {
	if in.AssetContractAddress == "" {
		return nil, codes.ErrInvalidArgument.Newf("invalid asset_contract_address params")
	}
	token_id := in.TokenId
	if token_id == "" {
		token_id = defaultTokenID
	}
	includeOrders := "false"
	if in.IncludeOrders == "true" {
		includeOrders = "true"
	}
	apiStr := "asset"
	queryParams := "?include_orders=" + includeOrders
	api := openseaApiUrl
	if in.Network == "test" {
		api = openseaTestApiUrl
	}
	apiUrl := api + path.Join(apiStr, in.AssetContractAddress, token_id, queryParams)
	out, err := doOpenseaApi(ctx, apiUrl, in.Network)
	if err != nil {
		return nil, codes.ErrTooManyRequests.Newf(err.Error())
	}
	return out, nil
}

func GetSingleAssetOut(ctx context.Context, in *SingleAssetInput) (*models.Asset, error) {
	res := &models.Asset{}
	out, err := GetSingleAssetApi(ctx, in)
	if err != nil {
		return nil, codes.ErrTooManyRequests.Newf(err.Error())
	}
	out.Restult(res)
	return res, nil
}

func GetSingleAsset(ctx context.Context, in *SingleAssetInput) (string, error) {
	out, err := GetSingleAssetApi(ctx, in)
	if err != nil {
		return "", codes.ErrTooManyRequests.Newf(err.Error())
	}
	return out.String(), nil
}

func doOpenseaApi(ctx context.Context, api string, network string) (*HttpResult, error) {
	/* proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:1087")
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport} */
	client := http.DefaultClient
	req, _ := http.NewRequest("GET", api, nil)

	if network != "test" {
		if xApiKey == "" {
			return nil, errors.New("invalid api")
		}
		req.Header.Add("X-API-KEY", xApiKey)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	httpResult := &HttpResult{}
	httpResult.status = res.Status
	httpResult.statusCode = res.StatusCode
	if res.StatusCode != http.StatusOK {
		logrus.Printf("url[%s],status: %s", api, res.Status)
		//return "", errors.New(res.Status)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	httpResult.body = body
	//json.Unmarshal(body, out)

	return httpResult, nil
}
