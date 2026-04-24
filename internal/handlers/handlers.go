package handlers

import (
	"chidis/internal/storage"
	"fmt"
	"strings"
)

var ErrQuit = fmt.Errorf("client wants to quit")

func PingHandler(args []string) (string, error) {
	return "+pong", nil
}
func EchoHandler(args []string) (string, error) {
	return "+" + strings.Join(args, " ") + "\n", nil
}
func QuitHandler(args []string) (string, error) {
	return "+OK\n", ErrQuit
}
func GetHandler(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("dont parameters")
	}
	if len(args) > 1 {
		return "", fmt.Errorf("many parameters (need 1)")
	}
	key := args[0]
	storage.Mutex.RLock()
	value, ok := storage.DataMap[key]
	storage.Mutex.RUnlock()
	if !ok {
		return "", fmt.Errorf("there is no value for this key")
	}
	return "+" + value + "\n", nil

}
func SetHandler(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("need 2 parameters")
	}

	storage.Mutex.Lock()
	storage.DataMap[args[0]] = args[1]
	storage.Mutex.Unlock()
	return "+OK\n", nil
}
