package storage

import (
	"os"
	"sync"
)

var IsLoading bool
var DataMap = make(map[string]interface{})
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
