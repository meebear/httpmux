package main

import (
//"net/http"
)

type params struct {
	vars map[string]string
}

type Router struct {
	roots map[string]*section
}

//type Handle func(w http.ResponseWriter, r *http.Request, pr *params)
type Handle func(s string)

func New() *Router {
	r := &Router{
		roots: make(map[string]*section),
	}
	return r
}

func (r *Router) Handle(method, path string, h Handle) error {
	rs, ok := r.roots[method]
	if !ok {
		rs, _ = newSection(nil, "/")
		r.roots[method] = rs
	}

	return rs.addRoute(path, h)
}
