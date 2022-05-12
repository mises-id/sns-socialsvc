package opensea_api

import (
	"context"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/sirupsen/logrus"
)

type (
	AssetModel struct {
		ID       uint64 `json:"id"`
		ImageUrl string `json:"image_url"`
		Name     string `json:"name"`
	}

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
	ListAssetInput struct {
		Owner   string
		Limit   uint64
		Cursor  string
		Network string
	}
	ListAssetOutput struct {
		Assets   []*AssetModel `json:"assets"`
		Next     string        `json:"next"`
		Previous string        `json:"previous"`
	}
)

var (
	openseaApiUrl     = "https://api.opensea.io/api/v1/"
	openseaTestApiUrl = "https://testnets-api.opensea.io/api/v1/"
	defaultTokenID    string
	xApiKey           string
)

func init() {
	defaultTokenID = "1"
	xApiKey = "a5c5d9c4d27f463e9baf74972266f666"
}

func ListAsset(ctx context.Context, in *ListAssetInput) (string, error) {
	if in.Owner == "" {
		return "", codes.ErrInvalidArgument.Newf("invalid owner params")
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
	out, err := doOpenseaApi(ctx, apiUrl, in.Network)
	if err != nil {
		return "", codes.ErrTooManyRequests.Newf(err.Error())
	}
	return out, nil
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
	return out, nil
}

func GetSingleAsset(ctx context.Context, in *SingleAssetInput) (string, error) {
	if in.AssetContractAddress == "" {
		return "", codes.ErrInvalidArgument.Newf("invalid asset_contract_address params")
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
		return "", codes.ErrTooManyRequests.Newf(err.Error())
	}
	return out, nil
}

func doOpenseaApi(ctx context.Context, api string, network string) (string, error) {
	/* proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:1087")
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport} */
	client := http.DefaultClient
	req, _ := http.NewRequest("GET", api, nil)

	if network != "test" {
		req.Header.Add("X-API-KEY", xApiKey)
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		logrus.Printf("url[%s],status: %s", api, res.Status)
		//return "", errors.New(res.Status)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	//json.Unmarshal(body, out)

	return string(body), nil
}
