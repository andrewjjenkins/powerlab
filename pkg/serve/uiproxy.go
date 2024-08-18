package serve

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (s *server) registerUiServer(r *gin.RouterGroup, sockNode *gin.RouterGroup) {
	u, err := url.Parse("http://localhost:3000/")
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	r.Any("/*any", func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	sockNode.Any("/*any", func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}
