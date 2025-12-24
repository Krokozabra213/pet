package jwtv1

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

func JWTFloatToTime(floatTime float64) time.Time {
	sec := int64(floatTime)
	nsec := int64((floatTime - float64(sec)) * 1e9)
	return time.Unix(sec, nsec).UTC()
}

func generateTokenID() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	tokenID := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
	return tokenID, nil
}
