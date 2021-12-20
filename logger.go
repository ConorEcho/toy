package toy

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *RequestContext) {
		t := time.Now()
		c.Next()

		log.Printf("[%d] %s cost %v", c.statusCode, c.Request.RequestURI, time.Since(t))
	}
}
