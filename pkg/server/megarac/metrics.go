package megarac

import (
	"fmt"
	"strings"
	"time"
)

func (api *Api) GetMetrics() (string, error) {
	sensorsInternal, err := api.getSensorsInternal()
	if err != nil {
		return "", err
	}

	timestamp := time.Now().UnixMilli()
	var builder strings.Builder

	for _, sensor := range *sensorsInternal {
		fmt.Fprintf(
			&builder,
			"powerlab_server{name=\"%s\", sensor=\"%s\"} %.4f %d\n",
			api.Name(),
			sensor.Name,
			sensor.Reading,
			timestamp,
		)
	}
	return builder.String(), nil
}
