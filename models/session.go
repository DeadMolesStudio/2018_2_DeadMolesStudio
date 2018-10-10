package models

import (
	"sync"
)

type Sessions struct {
	sync.Mutex
	Sessions map[string]int
}
