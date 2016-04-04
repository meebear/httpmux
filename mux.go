package httpmux

import (
	"net/http"
)

type Mux struct {
	roots map[string]*section

	NotFound http.Handler

	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, *Context)
}

type HandlerFunc func(w http.ResponseWriter, req *http.Request, ctx *Context)

func (hf HandlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request, ctx *Context) {
	hf(w, req, ctx)
}

func New() *Mux {
	m := &Mux{
		roots: make(map[string]*section),
	}
	return m
}

func (m *Mux) Get(path string, h Handler) {
	m.Handle("GET", path, h)
}

func (m *Mux) Post(path string, h Handler) {
	m.Handle("POST", path, h)
}

func (m *Mux) Head(path string, h Handler) {
	m.Handle("HEAD", path, h)
}

func (m *Mux) Options(path string, h Handler) {
	m.Handle("OPTIONS", path, h)
}

func (m *Mux) Put(path string, h Handler) {
	m.Handle("PUT", path, h)
}

func (m *Mux) Patch(path string, h Handler) {
	m.Handle("PATCH", path, h)
}

func (m *Mux) Delete(path string, h Handler) {
	m.Handle("DELETE", path, h)
}

func (m *Mux) Handle(method, path string, h Handler) {
	rs, _ := m.roots[method]
	if rs == nil {
		errmsg := ""
		rs, errmsg = newSection(nil, "/")
		if errmsg != "" {
			panic(errmsg)
		}
		m.roots[method] = rs
	}
	rs.addRoute(path, h)
}

/*
func (m *Mux) ServeFiles(path string, root string) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path " + path)
	}
	fileServer := http.FileServer(http.Dir(root))
	m.Get(path, func(w http.ResponseWriter, req *http.Request, ctx *Context) {
		req.URL.Path = ctx.ParamByName("filepath")
		fileServer.ServeHTTP(w, req)
	})
}
*/

func (m *Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if m.PanicHandler != nil {
		defer m.reCover(w, req)
	}

	path := req.URL.Path

	if rs := m.roots[req.Method]; rs != nil {
		ctx := Context{}
		h := rs.findRoute(path, &ctx)
		if h != nil {
			h.ServeHTTP(w, req, &ctx)
			return
		}
	}

	if m.NotFound != nil {
		m.NotFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func (m *Mux) reCover(w http.ResponseWriter, req *http.Request) {
	if rc := recover(); rc != nil {
		m.PanicHandler(w, req, rc)
	}
}
