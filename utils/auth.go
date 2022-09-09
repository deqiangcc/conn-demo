package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// 签名
func Sign(appid, secret string) string {
	return base64.URLEncoding.EncodeToString(GenHmac(appid, secret))
}

// Hmac加密
func GenHmac(appID, secret string) []byte {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(appID))
	return hash.Sum(nil)
}

// 验证签名
func SignVerify(appID, secret, sign string) bool {
	expectedSign := GenHmac(appID, secret)
	givenSign, err := base64.URLEncoding.DecodeString(sign)
	if err != nil {
		return false
	}

	return hmac.Equal(expectedSign, givenSign)
}
