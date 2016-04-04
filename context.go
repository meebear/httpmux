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

func (ctx *Context) ParamByIdx(idx int) string {
	if idx >= 0 && idx < len(ctx.params) {
		p := ctx.params[idx]
		if len(p.vars) > 0 {
			return p.vars[0]
		}
	}
	return ""
}

func (ctx *Context) ParamByName(name string) string {
	p, ok := ctx.paramMap[name]
	if ok {
		if len(p.vars) > 0 {
			return p.vars[0]
		}
	}
	return ""
}
