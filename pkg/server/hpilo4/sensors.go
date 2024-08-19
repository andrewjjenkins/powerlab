package hpilo4

import (
	"fmt"

	"github.com/andrewjjenkins/powerlab/pkg/model"
)

type TemperatureSensorResponse struct {
	Label          string `json:"label"`
	XPosition      int    `json:"xposition"`
	YPosition      int    `json:"yposition"`
	Location       string `json:"location"`
	Status         string `json:"status"`
	CurrentReading int    `json:"currentreading"`
	Caution        int    `json:"caution"`
	Critical       int    `json:"critical"`
	TempUnit       string `json:"temp_unit"`
}

type TemperatureSensorsResponse struct {
	HostPowerState string                      `json:"hostpwr_state"`
	Temperatures   []TemperatureSensorResponse `json:"temperature"`
}

func (api *Api) getSensorsInternal() (*TemperatureSensorsResponse, error) {
	var temps TemperatureSensorsResponse
	err := api.GetJson("/json/health_temperature", &temps)
	if err != nil {
		return nil, err
	}
	return &temps, nil
}

func (api *Api) GetSensors() (*model.ServerSensorReadings, error) {
	temps, err := api.getSensorsInternal()
	if err != nil {
		return nil, err
	}

	tempsByName := make(map[string]TemperatureSensorResponse)
	for _, t := range temps.Temperatures {
		if _, ok := tempsByName[t.Label]; ok {
			return nil, fmt.Errorf("duplicate temp sensor %s", t.Label)
		}
		tempsByName[t.Label] = t
	}

	for _, l := range []string{"02-CPU 1", "12-PS 2 Inlet"} {
		if _, ok := tempsByName[l]; !ok {
			return nil, fmt.Errorf("no temp sensor %s", l)
		}
	}

	return &model.ServerSensorReadings{
		CpuTemp:     float64(tempsByName["02-CPU 1"].CurrentReading),
		ChassisTemp: float64(tempsByName["12-PS 2 Inlet"].CurrentReading),
		FanSpeed:    0.0,
		PowerWatts:  0.0,
	}, nil
}

func (api *Api) GetSensorsRaw() (interface{}, error) {
	return api.getSensorsInternal()
}
