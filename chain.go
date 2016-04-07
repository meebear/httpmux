// Copyright (c) 2016 Alan Kang. All rights reserved.
// Chain of middlewares support for httpmux

package httpmux

import (
	"net/http"
)

type Chain struct {
	mws     []Handler // slice of middlewares
	h       Handler   // final handler of request
	nextIdx int
}

func NewChain() *Chain {
	return &Chain{}
}

func (c *Chain) PrependMiddlewares(mws ...Handler) *Chain {
	var m []Handler
	m = append(m, mws...)
	m = append(m, c.mws...)
	c.mws = m
	return c
}

func (c *Chain) AppendMiddlewares(mws ...Handler) *Chain {
	c.mws = append(c.mws, mws...)
	return c
}

func (c *Chain) Use(h Handler) *Chain {
	c.h = h
	return c
}

func (c *Chain) Next(w http.ResponseWriter, req *http.Request, ctx *Context) {
	var h Handler
	if c.nextIdx >= 0 && c.nextIdx < len(c.mws) {
		h = c.mws[c.nextIdx]
	} else if c.nextIdx == len(c.mws) {
		h = c.h
	} else {
		return
	}

	c.nextIdx++

	if h != nil {
		h.ServeHTTP(w, req, ctx)
	}
}

func (c *Chain) ServeHTTP(w http.ResponseWriter, req *http.Request, ctx *Context) {
	c.nextIdx = 0
	c.Next(w, req, ctx)
}
