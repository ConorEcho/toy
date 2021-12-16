package toy

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *RequestContext) {
		t := time.Now()
		c.Next()

		log.Printf(" %s in %v", c.Request.RequestURI, time.Since(t))
	}
}
