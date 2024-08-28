package main

import (
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type client struct {
	rwMux    sync.RWMutex
	userData map[int64]map[string]float64
}

func (c *client) getUserData(ctx *ext.Context, key string) (float64, bool) {
	c.rwMux.RLock()
	defer c.rwMux.RUnlock()

	if c.userData == nil {
		return 0, false
	}

	userData, ok := c.userData[ctx.EffectiveUser.Id]
	if !ok {
		return 0, false
	}

	v, ok := userData[key]
	return v, ok
}

func (c *client) setUserData(ctx *ext.Context, key string, val float64) {
	c.rwMux.Lock()
	defer c.rwMux.Unlock()

	if c.userData == nil {
		c.userData = map[int64]map[string]float64{}
	}

	_, ok := c.userData[ctx.EffectiveUser.Id]
	if !ok {
		c.userData[ctx.EffectiveUser.Id] = map[string]float64{}
	}
	c.userData[ctx.EffectiveUser.Id][key] = val
}
