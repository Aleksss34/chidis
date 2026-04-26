package storage

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var IsLoading bool
var DataMap = make(map[string]interface{})
var ExpiredMap = make(map[string]time.Time)
var Mutex = sync.RWMutex{}
var DBFile *os.File

func InitStorage() error {
	var err error
	DBFile, err = os.OpenFile("internal/data/db.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Save(setRequest string) error {

	if IsLoading {
		return nil
	}
	if _, err := DBFile.Write([]byte(setRequest)); err != nil {
		return err
	}
	return DBFile.Sync()

}
func TickExpired() {
	ticker := time.NewTicker(500 * time.Millisecond)

	for range ticker.C {
		Mutex.Lock()
		now := time.Now()
		count := 0
		limit := 100
		for key, exp := range ExpiredMap {
			if exp.Before(now) {
				count++
				if count > limit {
					break
				}
				log.Printf("[TICKER] Удален просроченный ключ: %s\n", key)
				delete(DataMap, key)
				delete(ExpiredMap, key)
				err := Save(fmt.Sprintf("DEL %s\n", key))
				if err != nil {
					log.Println("не удалось записать удаление от тикера в файл")
				}
			}

		}
		Mutex.Unlock()
	}
}
