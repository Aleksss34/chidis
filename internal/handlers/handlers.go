package handlers

import (
	"fmt"
	"strings"
)

var ErrQuit = fmt.Errorf("client wants to quit")

func PingHandler(args []string) (string, error) {
	return "+pong", nil
}
func EchoHandler(args []string) (string, error) {
	return "+" + strings.Join(args, " "), nil
}
func QuitHandler(args []string) (string, error) {
	return "+OK\n", ErrQuit
}
