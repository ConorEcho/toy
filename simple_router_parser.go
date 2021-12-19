package toy

import "strings"

type simpleMatcher struct {
	routeMap     map[string]struct{}
	vars         map[string]string
	matchedRoute string
}

func (m *simpleMatcher) GetMatchedVars() map[string]string {
	return m.vars
}

func (m *simpleMatcher) GetMatchedRoute() string {
	return m.matchedRoute
}

func NewSimpleParser() *simpleMatcher {
	return &simpleMatcher{
		routeMap: make(map[string]struct{}),
	}
}

func (m *simpleMatcher) Add(method string, route string) {
	m.routeMap[route] = struct{}{}
}

func (m *simpleMatcher) matchRoute(cur, target string) (map[string]string, bool) {
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

func (m *simpleMatcher) Match(method string, route string) (matched bool) {
	m.resetMatchedResult()

	for r, _ := range m.routeMap {
		if m.vars, matched = m.matchRoute(route, r); matched {
			m.matchedRoute = r
			return
		}
	}

	return
}

func (m *simpleMatcher) resetMatchedResult() {
	m.matchedRoute = ""
	m.vars = make(map[string]string)
}
