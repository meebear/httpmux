package main

import (
	"net/http"
)

type param struct {
	vars []string
}

type Context struct {
	params   []*param
	paramMap map[string]*param
	Vars     map[string]string
}

type Mux struct {
	roots map[string]*section
}

type Handle func(w http.ResponseWriter, r *http.Request, ctx *Context)

func New() *Mux {
	r := &Mux{
		roots: make(map[string]*section),
	}
	return r
}

func (r *Mux) Handle(method, path string, h Handle) error {
	rs, ok := r.roots[method]
	if !ok {
		rs, _ = newSection(nil, "/")
		r.roots[method] = rs
	}

	return rs.addRoute(path, h)
}

func (r *Mux) findRoute(method, path string, ctx *Context) (Handle, error) {
	return nil, nil
}

func (r *Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//ctx := Context{}
}
