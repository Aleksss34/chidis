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
	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("key is not a string")
	}
	return "+" + strValue + "\n", nil

}
func SetHandler(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("need 2 parameters")
	}
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	storage.DataMap[args[0]] = args[1]
	setReq := fmt.Sprintf("SET %s %s\n", args[0], args[1])
	if err := storage.Save(setReq); err != nil {
		return "", err
	}

	return "+OK\n", nil
}
func DelHandler(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("вы не ввели аргумент")
	}

	delReq := "DEL"
	count := 0
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	for _, arg := range args {
		if _, ok := storage.DataMap[arg]; ok {
			count++
			delete(storage.DataMap, arg)
			delReq += fmt.Sprintf(" %s", arg)
		}

	}
	if count != 0 {
		if err := storage.Save(delReq + "\r\n"); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf(":%d\r\n", count), nil
}
