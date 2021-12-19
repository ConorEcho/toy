package toy

import "net/http"

type routeMatcher interface {
	Add(method string, route string)
	Match(method string, route string) (matched bool)
	GetMatchedVars() map[string]string
	GetMatchedRoute() string

	resetMatchedResult()
}

type routeStorage interface {
	Store(method string, route string, handlers HandlerChain)
	GetHandlers(method string, route string) HandlerChain
}

type router struct {
	storage routeStorage
	matcher routeMatcher
}

func newRouter() *router {
	return &router{
		storage: NewHashStorage(),
		matcher: NewSimpleParser(),
	}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	r.matcher.Add(method, pattern)
	r.storage.Store(method, pattern, HandlerChain{handler})
}

func (r *router) handle(c *RequestContext) {
	method := c.GetMethod()
	path := c.GetPath()
	if matched := r.matcher.Match(method, path); matched {
		c.Params = r.matcher.GetMatchedVars()
		c.handlers = append(c.handlers, r.storage.GetHandlers(method, r.matcher.GetMatchedRoute())...)
	} else {
		c.handlers = append(c.handlers, func(c *RequestContext) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.GetPath())
		})
	}

	c.Next()
}
