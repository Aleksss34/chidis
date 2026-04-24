package storage

import "sync"

var DataMap = make(map[string]string)
var Mutex = sync.RWMutex{}
