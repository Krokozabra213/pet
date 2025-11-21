package tests

import "github.com/brianvoe/gofakeit/v6"

const (
	passDefaultLen = 10
)

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefaultLen)
}

func randomUsername() string {
	return gofakeit.Username()
}

func randomID(min, max int) int {
	return gofakeit.Number(min, max)
}
