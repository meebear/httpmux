package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type section struct {
	sName        string // section name without leading char if not raw type
	sType        sectionType
	regexp       string
	hasNonRawSub bool // only one non-raw sub section is allowed
	subs         map[string]*section
	ts           bool // trailing slash, useful if this is last section
	h            Handle
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
		if len(name) == 1 {
			return nil, "regexp empty"
		}
		if name[1] == '{' {
			if i := strings.Index(name, "}"); i == -1 {
				return nil, "regexp format error"
			} else {
				s.sName = name[2:i]
				s.regexp = name[i+1:]
			}
		} else {
			s.regexp = name[1:]
		}
	default:
		s.sType = SectionTypeRaw
		s.sName = name
	}

	if sParent != nil {
		fmt.Printf("pname=%s ptype=%s", sParent.sName, sParent.sType)
	}
	fmt.Printf(" name=%s sName=%s sType=%s regexp=%s\n", name, s.sName, s.sType, s.regexp)

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

func (rs *section) addRoute(path string, h Handle) error {
	if h == nil {
		return fmt.Errorf("handle not defined for path '%s'", path)
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
				return fmt.Errorf("error: addRoute: %s, %s", path, errmsg)
			}
			s.subs[p] = ss

		}
		s = ss
	}

	if s.h != nil {
		return fmt.Errorf("handle for '" + path + "' redefined")
	}
	s.h = h

	if s != rs {
		s.ts = strings.HasSuffix(path, "/")
	}

	return nil
}

func (rs *section) match(method, path string, ctx *Context) (Handle, error) {
	return nil, nil
}

func testHandle(w http.ResponseWriter, r *http.Request, ctx *Context) {
	fmt.Println("testHandle: ")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("path needed")
		os.Exit(1)
	}
	fmt.Println("add root")
	rs, _ := newSection(nil, "/")
	fmt.Println("add '/1/2/3'")
	err := rs.addRoute("/1/2/3", testHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("add '/1/2'")
	err = rs.addRoute("/1/2", testHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("add '/1/:nm'")
	err = rs.addRoute("/1/:nm", testHandle)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("add '/2/*wc'/8")
	err = rs.addRoute("/2/*wc/8", testHandle)
	if err != nil {
		fmt.Println(err)
	}
	err = rs.addRoute(os.Args[1], testHandle)
	if err != nil {
		fmt.Println(err)
	}
}
