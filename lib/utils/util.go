package utils

import (
	"math/rand"
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
