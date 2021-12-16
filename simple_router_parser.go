package toy

import "strings"

type simpleParser struct {
	routerMap map[string]struct{}
}

func NewSimpleParser() *simpleParser {
	return &simpleParser{make(map[string]struct{})}
}

func (p *simpleParser) insert(method string, route string) {
	p.routerMap[route] = struct{}{}
}

func (p *simpleParser) matchRoute(cur, target string) (map[string]string, bool) {
	cur = strings.Trim(cur, "/")
	target = strings.Trim(target, "/")

	params := make(map[string]string)

	if cur == target {
		return params, true
	}

	curSplit := strings.Split(cur, "/")
	targetSplit := strings.Split(target, "/")

	for i, part := range targetSplit {
		if i > len(curSplit) {
			return params, false
		}

		if part == "" {
			return params, false
		}

		if part[0] == ':' {
			params[part[1:]] = curSplit[i]
		} else if part[0] == '*' {
			params[part[1:]] = strings.Join(curSplit[i:], "/")
			return params, true
		} else if part != curSplit[i] {
			return nil, false
		}
	}

	if len(curSplit) != len(targetSplit) {
		return nil, false
	}

	return params, true
}

func (p *simpleParser) parse(method string, route string) (params map[string]string, matchRoute *string) {
	matched := false
	for r, _ := range p.routerMap {
		if params, matched = p.matchRoute(route, r); matched {
			matchRoute = &r
			return
		}
	}

	return
}
