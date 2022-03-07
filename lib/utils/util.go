package utils

import (
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func RandShuffle(slice []interface{}) {

	if len(slice) < 1 {
		return
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

}

func WirteLogDay(path, content string) error {
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
