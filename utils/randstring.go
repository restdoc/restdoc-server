package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var lowerCaseLetterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
var numberRunes = []rune("0123456789")

func LowerCaseRandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = lowerCaseLetterRunes[rand.Intn(len(lowerCaseLetterRunes))]
	}
	return string(b)
}

func RandNumberRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numberRunes[rand.Intn(len(numberRunes))]
	}
	return string(b)

}
