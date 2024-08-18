// Package serve Powerlab API
//
// Serves consolidated view of multiple IPMI/DRAC/iLO implementations.
//
//	Schemes: http, https
//	Host: localhost
//	BasePath: /api
//	Version: 1.0
//	Contact: Andrew Jenkins <andrewjjenkins@gmail.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//	- text/plain
//
// swagger:meta
package serve

import (
	"net/http"

	"github.com/gin-gonic/gin"

	serverModule "github.com/andrewjjenkins/powerlab/pkg/serve/server"
	serverServerModule "github.com/andrewjjenkins/powerlab/pkg/server"
)

//go:generate mkdir -p api
//go:generate go run github.com/go-swagger/go-swagger/cmd/swagger generate spec -m -o api/swagger.json

type server struct {
	e *gin.Engine
}

func Serve(inServer *http.Server, serverManager *serverServerModule.ServerManager) {
	s := &server{
		e: gin.Default(),
	}

	s.registerApiServer(s.e.Group("/api"))

	ss := serverModule.NewServer(serverManager)
	ss.RegisterApiServer(s.e.Group("/api"))

	s.registerUiServer(s.e.Group("/ui"), s.e.Group("/sockjs-node"))

	inServer.Handler = s.e
}

// VersionResponse describes the version of the server and API to clients
// swagger:model versionResponse
type VersionResponseBody struct {
	// Required: true
	Major int `json:"major"`
	// Required: true
	Minor int `json:"minor"`
	// Required: true
	Build string `json:"build"`
}

func (s *server) wrap(f func(s *server, c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(s, c)
	}
}

func (s *server) registerApiServer(r *gin.RouterGroup) {
	// swagger:route GET /api/version getVersion
	//
	// Returns the version of the server and API
	//
	//
	//   Produces:
	//   - application/json
	//
	//   Responses:
	//     200: body:versionResponse
	r.GET("/version", getVersion)
}

func getVersion(c *gin.Context) {
	c.JSON(http.StatusOK, &VersionResponseBody{1, 0, "xxx"})
}
