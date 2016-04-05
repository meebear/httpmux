// Copyright (c) 2016 Alan Kang. All rights reserved.
// Context stores URL parameters, params can be accessed by index
// or name if named.
// ctx.Vars can be used to store random data by handlers, won't be used
// by httpmux

package httpmux

type param struct {
	vars []string
}

type Context struct {
	Chain
	params   []*param
	paramMap map[string]*param
	Vars     map[string]string
}

func (ctx *Context) setParam(name, value string) {
	p := &param{}
	p.vars = append(p.vars, value)
	ctx.params = append(ctx.params, p)
	if len(name) > 0 {
		if ctx.paramMap == nil {
			ctx.paramMap = make(map[string]*param)
		}
		ctx.paramMap[name] = p
	}
}

func (ctx *Context) ParamsByIdx(idx int) []string {
	if idx >= 0 && idx < len(ctx.params) {
		p := ctx.params[idx]
		return p.vars
	}
	return nil
}

func (ctx *Context) ParamByIdx(idx int) string {
	vars := ctx.ParamsByIdx(idx)
	if len(vars) > 0 {
		return vars[0]
	}
	return ""
}

func (ctx *Context) ParamsByName(name string) []string {
	if p, ok := ctx.paramMap[name]; ok {
		return p.vars
	}
	return nil
}

func (ctx *Context) ParamByName(name string) string {
	vars := ctx.ParamsByName(name)
	if len(vars) > 0 {
		return vars[0]
	}
	return ""
}
