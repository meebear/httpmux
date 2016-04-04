package httpmux

import (
	"net/http"
)

type Chain struct {
	mws     []Handler // middlewares
	h       Handler
	nextIdx int
}

func NewChain() *Chain {
	return &Chain{}
}

func (c *Chain) appendMiddlewares(mws ...Handler) {
	c.mws = append(c.mws, mws...)
}

func (c *Chain) Use(h Handler) {
	c.h = h
}

func (c *Chain) Next(w http.ResponseWriter, req *http.Request, ctx *Context) {
	var h Handler
	if c.nextIdx < len(c.mws) {
		h = c.mws[c.nextIdx]
		c.nextIdx++
	} else {
		h = c.h
	}

	if h != nil {
		h.ServeHTTP(w, req, ctx)
	}
}

func (c *Chain) ServeHTTP(w http.ResponseWriter, req *http.Request, ctx *Context) {
	c.nextIdx = 0
	c.Next(w, req, ctx)
}
