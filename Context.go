package kube

import "sync"

type Context struct {
	*sync.Map
}

func (c *Context) Create(key, value interface{}) { c.Map.Store(key, value) }
func (c *Context) Update(key, value interface{}) { c.Map.Store(key, value) }
func (c *Context) Delete(key interface{})        { c.Map.Delete(key) }
func (c *Context) Get(key interface{}) interface{} {
	var val, _ = c.Map.Load(key)
	return val
}
