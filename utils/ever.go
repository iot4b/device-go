package utils

import (
	"device-go/dsm"
	"math/rand"
)

func RandomString(length int) string {
	var letters = []rune("1234567890abcdef")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateRandomAddress() dsm.EverAddress {
	return dsm.EverAddress("0:" + RandomString(64))
}
