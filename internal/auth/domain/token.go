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

type LogoutInput struct {
	refreshToken string
}

func NewLogoutInput(token string) *LogoutInput {
	return &LogoutInput{
		refreshToken: token,
	}
}

func (input *LogoutInput) GetRefreshToken() string {
	return input.refreshToken
}

type RefreshInput struct {
	refreshToken string
}

func NewRefreshInput(token string) *RefreshInput {
	return &RefreshInput{
		refreshToken: token,
	}
}

func (input *RefreshInput) GetRefreshToken() string {
	return input.refreshToken
}

type LoginOutput struct {
	accessT  string
	refreshT string
}

func NewLoginOutput(accessToken string, refreshToken string) *LoginOutput {
	return &LoginOutput{
		accessT:  accessToken,
		refreshT: refreshToken,
	}
}

func (input *LoginOutput) GetAccess() string {
	return input.accessT
}

func (input *LoginOutput) GetRefresh() string {
	return input.refreshT
}

type RefreshOutput struct {
	accessT  string
	refreshT string
}

func NewRefreshOutput(accessToken string, refreshToken string) *RefreshOutput {
	return &RefreshOutput{
		accessT:  accessToken,
		refreshT: refreshToken,
	}
}

func (input *RefreshOutput) GetAccess() string {
	return input.accessT
}

func (input *RefreshOutput) GetRefresh() string {
	return input.refreshT
}
