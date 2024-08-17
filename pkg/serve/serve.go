package serve

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mkdir -p api
//go:generate go run github.com/go-swagger/go-swagger/cmd/swagger generate spec -m -o api/swagger.json

type server struct {
	e *gin.Engine
}

func Serve(inServer *http.Server) {
	s := &server{
		e: gin.Default(),
	}

	s.registerApiServer(s.e.Group("/api"))

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
