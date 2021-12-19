package toy

import "net/http"

type routeMatcher interface {
	Add(method string, route string)
	Match(method string, route string) (matched bool)
	GetMatchedVars() map[string]string
	GetMatchedRoute() string

	resetMatchedResult()
}

type router struct {
	handlers map[string]HandlerFunc
	matcher  routeMatcher
}

func newRouter(matcher routeMatcher) *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		matcher:  matcher,
	}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	r.matcher.Add(method, pattern)
	r.handlers[pattern] = handler
}

func (r *router) handle(c *RequestContext) {
	method := c.GetMethod()
	path := c.GetPath()
	if matched := r.matcher.Match(method, path); matched {
		c.Params = r.matcher.GetMatchedVars()
		c.handlers = append(c.handlers, r.handlers[r.matcher.GetMatchedRoute()])
	} else {
		c.handlers = append(c.handlers, func(c *RequestContext) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.GetPath())
		})
	}

	c.Next()
}
