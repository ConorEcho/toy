package toy

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type RequestContext struct {
	Request        *http.Request
	responseWriter http.ResponseWriter

	Params map[string]string

	fullPath string

	handlers []HandlerFunc
	index    int
}

func NewContext(w http.ResponseWriter, req *http.Request, method string) *RequestContext {
	return &RequestContext{
		responseWriter: w,
		Request:        req,
		index:          -1,
		fullPath:       req.URL.Path,
	}
}

func (c *RequestContext) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *RequestContext) GetPath() string {
	return c.fullPath
}

func (c *RequestContext) GetMethod() string {
	return c.Request.Method
}

func (c *RequestContext) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *RequestContext) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *RequestContext) Status(code int) {
	c.responseWriter.WriteHeader(code)
}

func (c *RequestContext) SetHeader(key string, value string) {
	c.responseWriter.Header().Set(key, value)
}

func (c *RequestContext) JSON(statusCode int, body H) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(statusCode)
	encoder := json.NewEncoder(c.responseWriter)
	if err := encoder.Encode(body); err != nil {
		http.Error(c.responseWriter, err.Error(), 500)
	}
}

func (c *RequestContext) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.responseWriter.Write([]byte(html))
}

func (c *RequestContext) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.responseWriter.Write([]byte(fmt.Sprintf(format, values...)))
}
