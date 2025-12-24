package hmacv1

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type HMACHasher struct {
	secret []byte
}

func New(secret []byte) HMACHasher {
	return HMACHasher{
		secret: secret,
	}
}

func (n HMACHasher) HashJWTTokenHMAC(token string) string {
	h := hmac.New(sha256.New, n.secret)
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
