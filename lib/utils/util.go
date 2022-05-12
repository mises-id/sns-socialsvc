package utils

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	misesidPrefix    = "did:mises:"
	channelUrlPrefix = "ch_"
	ethEmptyAddress  = "0x0000000000000000000000000000000000000000"
)

func InArrayObject(elem primitive.ObjectID, arr []primitive.ObjectID) int {
	var index int
	index = -1
	for k, v := range arr {
		if v == elem {
			index = k
			break
		}
	}
	return index
}

func EthAddressToEIPAddress(address string) string {
	addrLowerStr := strings.ToLower(address)
	if strings.HasPrefix(addrLowerStr, "0x") {
		addrLowerStr = addrLowerStr[2:]
		address = address[2:]
	}
	var binaryStr string
	addrBytes := []byte(addrLowerStr)
	hash256 := crypto.Keccak256Hash([]byte(addrLowerStr))

	for i, e := range addrLowerStr {
		if e >= '0' && e <= '9' {
			continue
		} else {
			binaryStr = fmt.Sprintf("%08b", hash256[i/2])
			if binaryStr[4*(i%2)] == '1' {
				addrBytes[i] -= 32
			}
		}
	}

	return "0x" + string(addrBytes)
}

func EthAddressIsEmpty(eth_addresses string) bool {
	return eth_addresses == ethEmptyAddress
}

func UMisesToMises(umises uint64) (mises float64) {

	if umises == 0 {
		return mises
	}
	mises = float64(umises) / float64(1000000)
	return mises

}

func GetRand(min, max int) int {
	return int(rand.Int63n(int64(max-min))) + min
}

func AddMisesidProfix(misesid string) string {
	if misesid == "" {
		return misesid
	}
	if !strings.HasPrefix(misesid, misesidPrefix) {
		return misesidPrefix + misesid
	}
	return misesid
}
func RemoveMisesidProfix(misesid string) string {
	if strings.HasPrefix(misesid, misesidPrefix) {
		return strings.TrimPrefix(misesid, misesidPrefix)
	}
	return misesid
}
func AddChannelUrlProfix(channel_url string) string {
	if channel_url == "" {
		return channel_url
	}
	if !strings.HasPrefix(channel_url, channelUrlPrefix) {
		return channelUrlPrefix + channel_url
	}
	return channel_url
}
func RemoveChannelUrlProfix(channel_url string) string {
	if strings.HasPrefix(channel_url, channelUrlPrefix) {
		return strings.TrimPrefix(channel_url, channelUrlPrefix)
	}
	return channel_url
}

func RandShuffle(slice []interface{}) {

	if len(slice) < 1 {
		return
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

}

func WirteLogDay(path string) error {
	content := time.Now().String() + "\n"
	arr := strings.Split(path, "/")
	filePath := strings.Join(arr[:len(arr)-1], "/")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	fileObj, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer fileObj.Close()
	if _, err := io.WriteString(fileObj, content); err == nil {
		return err
	}
	return nil
}
