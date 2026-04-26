package handlers

import (
	"chidis/internal/storage"
	"fmt"
	"strconv"
	"strings"
)

func LPushHandler(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("недостаточно аргументов, нужно 2")
	}
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	value, ok := storage.DataMap[args[0]]
	if !ok {
		value = make([]string, 0)

	}
	valueArr, ok := value.([]string)

	if !ok {
		return "", fmt.Errorf("the entered key is not a list")
	}
	newArgs := append(args[1:], valueArr...)
	storage.DataMap[args[0]] = newArgs
	str := strings.Join(args, " ")
	lpushReq := fmt.Sprintf("LPUSH %s\n", str)
	if err := storage.Save(lpushReq); err != nil {
		return "", err
	}
	return fmt.Sprintf(":%d\r\n", len(newArgs)), nil
}
func RPushHandler(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("недостаточно аргументов, нужно 2")
	}
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	value, ok := storage.DataMap[args[0]]
	if !ok {
		value = make([]string, 0)

	}
	valueArr, ok := value.([]string)

	if !ok {
		return "", fmt.Errorf("the entered key is not a list")
	}
	newArgs := append(valueArr, args[1:]...)
	storage.DataMap[args[0]] = newArgs
	str := strings.Join(args, " ")
	rpushReq := fmt.Sprintf("RPUSH %s", str)
	if err := storage.Save(rpushReq); err != nil {
		return "", err
	}
	return fmt.Sprintf(":%d\r\n", len(newArgs)), nil
}
func LRangeHandler(args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("need 3 arguments")
	}
	key := args[0]
	start, err := strconv.Atoi(args[1])
	if err != nil {
		return "", fmt.Errorf("start - integer number")
	}
	finish, err := strconv.Atoi(args[2])
	if err != nil {
		return "", fmt.Errorf("finish - integer number")
	}
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	value, ok := storage.DataMap[key]
	if !ok {
		return "*0\r\n", nil
	}
	valueArr, ok := value.([]string)
	if !ok {
		return "", fmt.Errorf("WRONGTYPE")
	}
	n := len(valueArr)
	if start < 0 {
		start += n
	}
	if finish < 0 {
		finish += n
	}
	if start < 0 {
		start = 0
	}
	if finish >= n {
		finish = n - 1
	}
	if start >= n || start > finish {
		return "*0\r\n", nil
	}
	count := finish - start + 1
	var str strings.Builder
	str.WriteString(fmt.Sprintf("*%d\r\n", count))
	for i := start; i <= finish; i++ {
		str.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(valueArr[i]), valueArr[i]))
	}
	return str.String(), nil
}

func LLenHandler(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("need 1 arguments")
	}

	_, ok := storage.DataMap[args[0]]
	if !ok {
		return "", fmt.Errorf("invalid key")
	}
	value, ok := storage.DataMap[args[0]].([]string)
	if !ok {
		return "", fmt.Errorf("value is not list")
	}

	return strconv.Itoa(len(value)), nil
}

func LPopHandler(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("need 1 arguments")
	}
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	_, ok := storage.DataMap[args[0]]
	if !ok {
		return "$-1\r\n", nil
	}
	value, ok := storage.DataMap[args[0]].([]string)
	if !ok {
		return "", fmt.Errorf("value is not list")
	}
	if len(value[1:]) == 0 {
		delete(storage.DataMap, args[0])

	} else {
		storage.DataMap[args[0]] = value[1:]
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(value[0]), value[0]), nil
}
func RPopHandler(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("need 1 arguments")
	}
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	_, ok := storage.DataMap[args[0]]
	if !ok {
		return "$-1\r\n", nil
	}
	value, ok := storage.DataMap[args[0]].([]string)
	if !ok {
		return "", fmt.Errorf("value is not list")
	}
	newValue := value[:len(value)-1]
	if len(newValue) == 0 {
		delete(storage.DataMap, args[0])

	} else {
		storage.DataMap[args[0]] = newValue
	}
	if err := storage.Save(fmt.Sprintf("RPOP %s\n", args[0])); err != nil {
		return "", fmt.Errorf("не удалос сохранить команду в файл")
	}

	return fmt.Sprintf("$%d\r\n%s\r\n", len(value[len(value)-1]), value[len(value)-1]), nil
}
