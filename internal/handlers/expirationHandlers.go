package handlers

import (
	"chidis/internal/storage"
	"fmt"
	"strconv"
	"time"
)

func ExpireHandler(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("вы не ввели аргумент")
	}
	if args[0] == "resetly" {
		exp, _ := strconv.ParseInt(args[2], 10, 64)
		storage.Mutex.Lock()
		storage.ExpiredMap[args[1]] = time.Unix(exp, 0)
		storage.Mutex.Unlock()
		return "+OK\r\n", nil
	}
	if len(args) != 2 {
		return "", fmt.Errorf("need 2 parameters")
	}
	if _, ok := storage.DataMap[args[0]]; !ok {
		return ":0\r\n", nil
	}
	t, err := strconv.Atoi(args[1])
	if err != nil {
		return "", err
	}
	exitedAt := time.Now().Add(time.Duration(t) * time.Second)
	storage.ExpiredMap[args[0]] = exitedAt
	err = storage.Save(fmt.Sprintf("Expire resetly %s %d\n", args[0], exitedAt.Unix()))
	if err != nil {
		return "", err
	}
	return ":1\r\n", nil
}
func TtlHandler(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("need 1 parameters")
	}
	expiredAs, ok := storage.ExpiredMap[args[0]]
	if !ok {
		return ":-1\r\n", nil
	}
	resp := expiredAs.Unix() - time.Now().Unix()
	if resp <= 0 {
		return ":-1\r\n", nil
	}
	return fmt.Sprintf(":%d\r\n", resp), nil

}
