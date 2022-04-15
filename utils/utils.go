package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"math/big"
	"net"
	"strings"

	Models "restdoc-models/models"
	"restdoc/config"
)

const DefaultColor = int32(4358125)

func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func GravatarUrl(email string) string {
	address := strings.TrimSpace(email)
	address = strings.ToLower(address)
	data := []byte(address)
	hash := md5.Sum(data)
	s := fmt.Sprintf("%x", hash)
	return fmt.Sprintf("%s/%s", config.DefaultConfig.GravatarUrl, s)
}

func GetHashedPassword(userId int64, password string) string {
	saltKey := fmt.Sprintf("%s%d", password, userId)
	salt := sha256.Sum256([]byte(saltKey))
	hashedSalt := fmt.Sprintf("%x", salt)
	sum := sha256.Sum256([]byte(hashedSalt + password))
	hashed := fmt.Sprintf("%x", sum)
	return hashed
}

func FormatColor(color int32) string {
	/*
		if color <= 0 {
			return fmt.Sprintf("#%06x", DefaultColor)
		}
	*/
	return fmt.Sprintf("#%06x", color)
}

func FormatMethod(methodType int16) string {
	method := "GET"
	switch methodType {
	case Models.METHOD_GET:
		method = "GET"
	case Models.METHOD_POST:
		method = "POST"
	case Models.METHOD_OPTION:
		method = "OPTION"
	default:
	}
	return method
}
