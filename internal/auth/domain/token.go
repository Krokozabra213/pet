package domain

import "time"

type BlackToken struct {
	Exp   time.Time
	Token string `gorm:"uniqueIndex;size:64"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
