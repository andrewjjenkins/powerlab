package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/andrewjjenkins/powerlab/pkg/server"
)

type ServerServer struct {
	manager *server.ServerManager
}

type ServerResponse struct {
	Name        string  `json:"name"`
	PowerStatus int     `json:"power_status"`
	PowerWatts  float64 `json:"power_watts"`
}

// ServersResponse describes the servers
// swagger:model serversResponse
type ServersResponse []ServerResponse

func NewServer(manager *server.ServerManager) ServerServer {
	return ServerServer{
		manager: manager,
	}
}

func (ss *ServerServer) wrap(f func(ss *ServerServer, c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(ss, c)
	}
}

func (ss *ServerServer) RegisterApiServer(r *gin.RouterGroup) {
	// swagger:route GET /api/servers getServers
	//
	// Returns the list of servers
	//
	//   Produces:
	//   - application/json
	//
	//   Response:
	//     200: body:serversResponse
	r.GET("/servers", ss.wrap(GetServers))
}

func GetServers(ss *ServerServer, c *gin.Context) {
	res := ServersResponse{}

	for _, s := range ss.manager.Servers {
		res = append(res, ServerResponse{
			Name:        s.Name(),
			PowerStatus: 1,
			PowerWatts:  251.7,
		})
	}
	c.JSON(http.StatusOK, res)
}
