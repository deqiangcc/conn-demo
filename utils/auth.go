package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func GenHmacStr(secret string) string {
	return base64.URLEncoding.EncodeToString(GenHmac(secret))
}

func GenHmac(secret string) []byte {
	hash := hmac.New(sha256.New, []byte(secret))
	//hash.Write([]byte(APP_ID))
	return hash.Sum(nil)
}

func HmacVerify(secret string, hmacStr string) bool {
	givenHmac, err := base64.URLEncoding.DecodeString(hmacStr)
	if err != nil {
		fmt.Println("decode hmac err:", err)
		return false
	}
	expectedHmac := GenHmac(secret)
	return hmac.Equal(expectedHmac, givenHmac)
}
