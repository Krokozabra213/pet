package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HashJWTTokenHMAC(token string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyJWTTokenHMAC(token, expectedHash string, secret []byte) bool {
	actualHash := HashJWTTokenHMAC(token, secret)
	return hmac.Equal([]byte(actualHash), []byte(expectedHash))
}
