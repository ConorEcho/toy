package toy

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*RequestContext)

type HandlerChain []HandlerFunc

type routerGroup struct {
	prefix   string
	parent   *routerGroup
	children []*routerGroup
	handlers []HandlerFunc
	engine   *Engine
}

// Engine implement the interface of ServeHTTP
type Engine struct {
	*routerGroup
	router *router
}

func New() *Engine {
	e := &Engine{
		router: newRouter(NewSimpleParser()),
		routerGroup: &routerGroup{
			children: []*routerGroup{},
		},
	}

	e.engine = e
	return e
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())

	return engine
}

func (group *routerGroup) Use(handler ...HandlerFunc) {
	group.handlers = append(group.handlers, handler...)
}

func (group *routerGroup) Group(prefix string) *routerGroup {
	g := &routerGroup{
		prefix:   group.prefix + prefix,
		parent:   group,
		children: []*routerGroup{},
		engine:   group.engine,
	}

	group.children = append(group.children, g)
	g.parent = group
	return g
}

func (group *routerGroup) addRoute(method, pattern string, handler HandlerFunc) {
	group.engine.router.addRoute(method, group.prefix+pattern, handler)
}

// GET defines the method to add GET request
func (group *routerGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *routerGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (group *routerGroup) Run(addr string) (err error) {
	return http.ListenAndServe(addr, group.engine)
}

func (group *routerGroup) getHandlers(path string) []HandlerFunc {
	h := group.handlers

	children := group.children
	for len(children) > 0 {
		t := []*routerGroup{}
		for _, child := range children {
			if strings.HasPrefix(path, child.prefix) {
				h = append(h, child.handlers...)
				t = child.children
				break
			}
		}

		children = t
	}

	return h
}

func (group *routerGroup) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := NewContext(w, req, req.Method)
	ctx.handlers = group.getHandlers(req.URL.Path)
	group.engine.router.handle(ctx)
}
