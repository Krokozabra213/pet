package domain

type User struct {
	ID       int64
	Username string
	PassHash []byte
}
