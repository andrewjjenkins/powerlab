package responsecache

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type Cache struct {
	Responses *expirable.LRU[string, interface{}]
}

func New() *Cache {
	c := Cache{
		Responses: expirable.NewLRU[string, interface{}](500, nil, time.Second*30),
	}
	return &c
}
