package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func getMetricsInternal(ss *ServerServer) (string, error) {
	var builder strings.Builder

	fmt.Fprintf(
		&builder,
		"# HELP powerlab_server A metric from a server IPMI\n"+
			"# TYPE powerlab_server gauge\n",
	)

	for _, server := range ss.manager.Servers {
		serverMetrics, err := server.GetMetrics()
		if err != nil {
			return "", err
		}
		builder.WriteString(serverMetrics)
	}
	return builder.String(), nil
}

func GetMetrics(ss *ServerServer, c *gin.Context) {
	metrics, err := getMetricsInternal(ss)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed getting metrics: %v", err))
	}

	c.Writer.Header().Set("content-type", "text/plain; version=0.0.4")
	c.Writer.Header().Set("content-length", fmt.Sprint(len(metrics)))
	bytes, err := c.Writer.Write([]byte(metrics))
	if err != nil {
		glog.Error("failed to write body: %v", err)
		return
	}
	if bytes != len(metrics) {
		glog.Error("failed to write entire body (%d < %d)", bytes, len(metrics))
		return
	}
}
