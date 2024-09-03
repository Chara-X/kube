package kube

import "sync"

type Context struct {
	*sync.Map
}

func (c *Context) Create(key string, value interface{}) { c.Map.Store(key, value) }
func (c *Context) Update(key string, value interface{}) { c.Map.Store(key, value) }
func (c *Context) Delete(key string)                    { c.Map.Delete(key) }
func (c *Context) Get(key string) interface{} {
	var val, _ = c.Map.Load(key)
	return val
}
