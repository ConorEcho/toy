package toy

import "net/http"

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

type router struct {
	handlers map[string]HandlerFunc
	parser   routerParser
}

type routerParser interface {
	insert(method string, route string)
	parse(method string, route string) (params map[string]string, matchRoute *string)
}

func newRouter(parser routerParser) *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		parser:   parser,
	}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	r.parser.insert(method, pattern)
	r.handlers[pattern] = handler
}

func (r *router) handle(c *Context) {
	method := c.GetMethod()
	path := c.GetPath()
	if params, route := r.parser.parse(method, path); route != nil {
		c.param = params
		c.handlers = append(c.handlers, r.handlers[*route])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.GetPath())
		})
	}

	c.Next()
}
