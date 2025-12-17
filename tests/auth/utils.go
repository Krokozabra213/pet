package tests

import "github.com/brianvoe/gofakeit/v6"

const (
	passDefaultLen = 10
)

func randomFakePassword() string {
	left := gofakeit.Password(true, true, false, false, false, 1)
	right := gofakeit.Password(true, true, true, false, false, passDefaultLen)
	return left + right
}

func randomUsername() string {
	return gofakeit.Username()
}

func randomID(min, max int) int {
	return gofakeit.Number(min, max)
}
