package workflow

import (
	"sync"

	"github.com/hpifu/go-kit/hconf"
)

type Ctx struct {
	mutex sync.RWMutex
	store *hconf.InterfaceStorage
}

func NewCtx() *Ctx {
	return &Ctx{
		store: hconf.NewInterfaceStorage(map[string]interface{}{}),
	}
}

func (c *Ctx) Get(key string) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.store.Get(key)
}

func (c *Ctx) Set(key string, val interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.store.Set(key, val)
}
