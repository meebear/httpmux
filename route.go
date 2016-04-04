package httpmux

import (
	"regexp"
	"strings"
)

type section struct {
	sName        string // section name without leading char if not raw type
	sType        sectionType
	regexp       *regexp.Regexp
	hasNonRawSub bool // only one non-raw sub section is allowed
	subs         map[string]*section
	ts           bool // trailing slash, useful if this is last section
	h            Handler
}

type flags uint32

type sectionType int

const (
	RedirectTrailingSlash = 1 << iota
)

const (
	SectionTypeRaw sectionType = iota
	SectionTypeWildCard
	SectionTypeMatch
	SectionTypeRegexp
)

func (s sectionType) String() string {
	n := "SectionTypeUnknown"
	switch s {
	case SectionTypeRaw:
		n = "SectionTypeRaw"
	case SectionTypeWildCard:
		n = "SectionTypeWildCard"
	case SectionTypeMatch:
		n = "SectionTypeMatch"
	case SectionTypeRegexp:
		n = "SectionTypeRexexp"
	}
	return n
}

func newSection(sParent *section, name string) (*section, string) {
	s := &section{}
	switch name[0] {
	case ':':
		s.sType = SectionTypeMatch
		s.sName = name[1:]
	case '*':
		s.sType = SectionTypeWildCard
		s.sName = name[1:]
	case '#':
		s.sType = SectionTypeRegexp
		var re string
		if len(name) == 1 {
			return nil, "regexp empty"
		}
		if name[1] == '{' {
			if i := strings.Index(name, "}"); i == -1 {
				return nil, "regexp format error"
			} else {
				s.sName = name[2:i]
				re = name[i+1:]
			}
		} else {
			re = name[1:]
		}
		if len(re) == 0 {
			return nil, "regexp empty"
		}
		var err error
		s.regexp, err = regexp.Compile(re)
		if err != nil {
			return nil, "regexp compile error"
		}
	default:
		s.sType = SectionTypeRaw
		s.sName = name
	}

	/*
		if sParent != nil {
			fmt.Printf("pname=%s ptype=%s", sParent.sName, sParent.sType)
		}
		fmt.Printf(" name=%s sName=%s sType=%s regexp=%s\n", name, s.sName, s.sType, s.regexp)
	*/

	if sParent != nil {
		if sParent.sType == SectionTypeWildCard {
			return nil, "wildcard not the last section"
		}
		if s.sType != SectionTypeRaw {
			if sParent.hasNonRawSub {
				return nil, "multiple non raw section"
			}
			sParent.hasNonRawSub = true
		}
	}

	return s, ""
}

func (rs *section) addRoute(path string, h Handler) {
	if h == nil {
		panic("handle not defined for path " + path)
	}
	if len(path) == 0 || path[0] != '/' {
		panic("path must begin with '/'")
	}

	s := rs
	ps := strings.Split(path, "/")
	for _, p := range ps {
		if len(p) == 0 {
			continue
		}
		if s.subs == nil {
			s.subs = make(map[string]*section)
		}

		p = strings.ToLower(p)
		ss, ok := s.subs[p]
		if !ok {
			errmsg := ""
			if ss, errmsg = newSection(s, p); errmsg != "" {
				panic("error: addRoute: " + path + " " + errmsg)
			}
			s.subs[p] = ss

		}
		s = ss
	}

	if s.h != nil {
		panic("handle for '" + path + "' redefined")
	}
	s.h = h

	if s != rs {
		s.ts = strings.HasSuffix(path, "/")
	}
}

func (s *section) match(ps []string, ctx *Context) (m bool, h Handler, stop bool) {
	switch s.sType {
	case SectionTypeWildCard:
		ctx.setParam(s.sName, strings.Join(ps, "/"))
		m, h, stop = true, s.h, true
	case SectionTypeMatch:
		ctx.setParam(s.sName, ps[0])
		m, h = true, s.h
	case SectionTypeRegexp:
		if s.regexp.Match([]byte(ps[0])) {
			ctx.setParam(s.sName, ps[0])
			m, h = true, s.h
		}
	case SectionTypeRaw:
		fallthrough
	default:
	}
	return
}

func (rs *section) findRoute(path string, ctx *Context) Handler {
	var h Handler
	isRoot := true
	s := rs
	ps := strings.Split(path, "/")
loop:
	for i, p := range ps {
		if len(p) == 0 {
			continue
		}
		isRoot = false
		if s.subs == nil {
			return nil
		}
		if ss, ok := s.subs[p]; ok { // matches raw
			h = ss.h
			s = ss
			continue
		}

		match, stop := false, false
		for _, ss := range s.subs {
			match, h, stop = ss.match(ps[i:], ctx)
			if match {
				if stop {
					break loop
				} else {
					s = ss
					continue loop
				}
			}
		}

		return nil
	}

	if isRoot {
		return s.h
	}
	return h
}
