package broker

import (
	"strconv"
	"strings"
)

type ClientDeps struct {
	Cli  Client
	Name string
}

func (deps *ClientDeps) Close() {
	deps.Cli.Close()
}

func NewClientDeps(name string) *ClientDeps {
	return &ClientDeps{Name: name, Cli: NewClient()}
}

type Users map[string]Client

func (u Users) Online() string {
	online := len(u)
	return strconv.Itoa(online)
}

func (u Users) String() string {
	var builder strings.Builder
	builder.WriteString("Users: ")
	for name := range u {
		builder.WriteString(name)
		builder.WriteString(" ")
	}
	builder.WriteRune('\n')
	return builder.String()
}
