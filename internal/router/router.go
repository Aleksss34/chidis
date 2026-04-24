package router

import (
	"chidis/internal/handlers"
	"fmt"
	"strings"
)

type Handler func(msgArgs []string) (string, error)

var HandlersMap = map[string]Handler{
	"PING": handlers.PingHandler,
	"ECHO": handlers.EchoHandler,
	"QUIT": handlers.QuitHandler,
	"GET":  handlers.GetHandler,
	"SET":  handlers.SetHandler,
}

func Routing(msg string) (string, error) {

	msg = strings.TrimSpace(msg)
	wordsMsg := strings.Fields(msg)
	if len(wordsMsg) == 0 {
		return "", fmt.Errorf("empty string")
	}
	cmd := strings.ToUpper(wordsMsg[0])

	handler, ok := HandlersMap[cmd]
	if !ok {
		return "", fmt.Errorf("error command")
	}
	return handler(wordsMsg[1:])
}
