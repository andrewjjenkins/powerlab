package hpilo4

import (
	"fmt"
	"strings"
	"time"
)

func (api *Api) GetMetrics() (string, error) {
	temps, err := api.getTemperaturesInternal()
	if err != nil {
		return "", err
	}
	powers, err := api.getPowerInternal()
	if err != nil {
		return "", err
	}

	timestamp := time.Now().UnixMilli()
	var builder strings.Builder

	for _, temp := range temps.Temperatures {
		fmt.Fprintf(
			&builder,
			"powerlab_server{name=\"%s\", sensor=\"%s\"} %.4f %d\n",
			api.Name(),
			temp.Label,
			float64(temp.CurrentReading),
			timestamp,
		)
	}

	fmt.Fprintf(
		&builder,
		// SYS_POWER matches the sensor name in megarac
		"powerlab_server{name=\"%s\", sensor=\"SYS_POWER\"} %.4f %d\n",
		api.Name(),
		float64(powers.PresentPowerReading),
		timestamp,
	)

	return builder.String(), nil
}
