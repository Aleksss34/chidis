package handlers

import (
	"chidis/internal/storage"
	"fmt"
	"strconv"
	"strings"
	"time"
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
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	timeKill, ok := storage.ExpiredMap[key]
	if ok && timeKill.Before(time.Now()) {
		delete(storage.ExpiredMap, key)
		delete(storage.DataMap, key)
		if err := storage.Save(fmt.Sprintf("DEL %s\n", key)); err != nil {
			return "", err
		}
		return "$-1\r\n", nil
	}
	value, ok := storage.DataMap[key]
	if !ok {
		return "$-1\r\n", nil
	}
	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("key is not a string")
	}

	return "$" + strValue + "\r\n", nil

}

func SetHandler(args []string) (string, error) {
	if args[0] == "resetly" {
		exp, _ := strconv.ParseInt(args[4], 10, 64)
		storage.Mutex.Lock()
		storage.ExpiredMap[args[1]] = time.Unix(exp, 0)
		storage.Mutex.Unlock()
		return "+OK\r\n", nil
	}
	if len(args) != 2 && len(args) != 4 {
		return "", fmt.Errorf("need 2 or 4 parameters")
	}
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	storage.DataMap[args[0]] = args[1]
	var setReq string
	if len(args) == 4 && strings.ToUpper(args[2]) == "EX" {
		t, err := strconv.Atoi(args[3])
		if err != nil {
			return "", err
		}
		storage.ExpiredMap[args[0]] = time.Now().Add(time.Duration(t) * time.Second)
		setReq = fmt.Sprintf("SET resetly %s %s %s %s\n", args[0], args[1], args[2], args[3])
	} else {
		setReq = fmt.Sprintf("SET %s %s\n", args[0], args[1])
	}

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
			delete(storage.ExpiredMap, arg)
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
