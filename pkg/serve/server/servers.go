package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/andrewjjenkins/powerlab/pkg/model"
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
// swagger:model ServersResponse
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

// SensorsResponse is sensor data (generic to all servers)
// swagger:model SensorsResponse
type SensorsResponse model.ServerSensorReadings

// SensorsRawResponse is the raw JSON of sensor data (server specific)
// swagger:model SensorsRawResponse
type SensorsRawResponse interface{}

func (ss *ServerServer) RegisterApiServer(r *gin.RouterGroup) {
	// swagger:route GET /api/servers getServers
	//
	// Returns the list of servers
	//
	//   Produces:
	//   - application/json
	//
	//   Responses:
	//     200: body:ServersResponse
	r.GET("/servers", ss.wrap(GetServers))

	// swagger:route GET /api/server/{name} getServerByName
	//
	// Returns a server
	//
	//   Produces:
	//   - application/json
	//
	//   Responses:
	//     200: body:ServerResponse
	r.GET("/servers/:name", ss.wrap(GetServerByName))

	// swagger:route GET /api/server/{name}/sensors getServerSensors
	//
	// Returns the sensors from a server
	//
	//   Produces:
	//   - application/json
	//
	//   Responses:
	//     200: body:SensorsResponse
	r.GET("/servers/:name/sensors", ss.wrap(GetServerSensors))

	// swagger:route GET /api/server/{name}/sensorsRaw getServerSensorsRaw
	//
	// Returns the sensors from a server as raw JSON
	//
	//   Produces:
	//   - application/json
	//
	//   Responses:
	//     200: body:SensorsRawResponse
	r.GET("/servers/:name/sensorsRaw", ss.wrap(GetServerSensorsRaw))

	// swagger:route GET /api/metrics getMetrics
	//
	// Returns the metrics from a server in prometheus format
	//
	//   Produces:
	//   - text/plain
	//
	//   Responses:
	//     200: body:string
	r.GET("/metrics", ss.wrap(GetMetrics))
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

func GetServerByName(ss *ServerServer, c *gin.Context) {
	name := c.Param("name")
	server, ok := ss.manager.Servers[name]
	if !ok {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("no server %s", name))
		return
	}
	c.JSON(http.StatusOK, ServerResponse{
		Name:        server.Name(),
		PowerStatus: 1,
		PowerWatts:  313.2,
	})
}
