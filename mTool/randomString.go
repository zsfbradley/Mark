package mTool

import (
	"math/rand"
	"time"
)

var (
	originSlice = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetRandString(size int) string {
	rs := ""
	for i := 0; i < size; i++ {
		c := originSlice[rand.Intn(len(originSlice))]
		rs += string(c)
	}
	return rs
}

func RandIntn(n int) int {
	return rand.Intn(n)
}
