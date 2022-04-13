package utils

import (
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	misesidPrefix    = "did:mises:"
	channelUrlPrefix = "ch_"
)

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
