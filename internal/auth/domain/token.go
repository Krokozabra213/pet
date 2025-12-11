package domain

import "time"

const TokenEntity = "Token"

type BlackToken struct {
	Exp   time.Time
	Token string `gorm:"uniqueIndex;size:64"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewBlackToken(exp time.Time, token string) *BlackToken {
	return &BlackToken{
		Exp:   exp,
		Token: token,
	}
}

func NewTokenPair(access, refresh string) *TokenPair {
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}
