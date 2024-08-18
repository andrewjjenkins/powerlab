package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetServerSensors(ss *ServerServer, c *gin.Context) {
	name := c.Param("name")
	server, ok := ss.manager.Servers[name]
	if !ok {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("no server %s", name))
		return
	}
	sensorOut, err := server.GetSensors()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting sensor data: %v", err))
	}
	c.JSON(http.StatusOK, sensorOut)
}

func GetServerSensorsRaw(ss *ServerServer, c *gin.Context) {
	name := c.Param("name")
	server, ok := ss.manager.Servers[name]
	if !ok {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("no server %s", name))
		return
	}
	rawSensorOut, err := server.GetSensorsRaw()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting sensor data: %v", err))
	}
	c.JSON(http.StatusOK, rawSensorOut)
}
