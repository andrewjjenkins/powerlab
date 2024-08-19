package hpilo4

import "github.com/andrewjjenkins/powerlab/pkg/model"

func (api *Api) GetSensors() (*model.ServerSensorReadings, error) {
	return &model.ServerSensorReadings{}, nil
}

func (api *Api) GetSensorsRaw() (interface{}, error) {
	return nil, nil
}
