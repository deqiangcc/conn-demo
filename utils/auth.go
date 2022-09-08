package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func GenToken(appid, secret string) string {
	return base64.URLEncoding.EncodeToString(GenHmac(appid, secret))
}

func GenHmac(appID, secret string) []byte {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(appID))
	return hash.Sum(nil)
}

func TokenVerify(appID, secret, token string) bool {
	expectedHmac := GenHmac(appID, secret)
	givenHmac, err := base64.URLEncoding.DecodeString(token)

	if err != nil {
		return false
	}

	return hmac.Equal(givenHmac, expectedHmac)
}
