package vm

import "fmt"

type Context struct {
	enclosing *Context
	values    map[string]any
}

func NewContext(enclosing *Context) *Context {
	c := new(Context)
	c.values = make(map[string]any)
	c.enclosing = enclosing
	return c
}

func (ctx *Context) Get(key string) (any, error) {
	val, ok := ctx.values[key]
	if ok {
		return val, nil
	}
	if ctx.enclosing != nil {
		return ctx.enclosing.Get(key)
	}
	return nil, fmt.Errorf("undefined variable '%s'", key)
}

func (ctx *Context) Assgin(key string, value any) error {
	_, flag := ctx.values[key]
	if flag {
		ctx.values[key] = value
		return nil
	}
	if ctx.enclosing != ctx {
		return ctx.enclosing.Assgin(key, value)
	}
	return fmt.Errorf("undefined variable '%s' when assign", key)
}

func (ctx *Context) Define(key string, value any) {
	ctx.values[key] = value
}
