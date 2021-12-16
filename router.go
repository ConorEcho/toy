package toy

import "net/http"

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

func (r *router) handle(c *RequestContext) {
	method := c.GetMethod()
	path := c.GetPath()
	if params, route := r.parser.parse(method, path); route != nil {
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[*route])
	} else {
		c.handlers = append(c.handlers, func(c *RequestContext) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.GetPath())
		})
	}

	c.Next()
}
