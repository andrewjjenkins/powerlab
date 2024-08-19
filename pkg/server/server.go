package server

import (
	"fmt"

	"github.com/andrewjjenkins/powerlab/pkg/model"
	"github.com/andrewjjenkins/powerlab/pkg/server/hpilo4"
	"github.com/andrewjjenkins/powerlab/pkg/server/megarac"
)

const (
	// These are from megarac, will need abstractions for other IPMI
	PowerCommandHardOff   = 0
	PowerCommandOn        = 1
	PowerCommandCycle     = 2
	PowerCommandHardReset = 3
	PowerCommandSoftOff   = 5
)

type Server interface {
	Login(username, password string) error
	Logout() error
	PowerCommand(command int) error
	GetSensorsRaw() (interface{}, error)
	GetSensors() (*model.ServerSensorReadings, error)
	GetMetrics() (string, error)
	Name() string
}

type ServerManager struct {
	Servers map[string]Server
}

func NewServerManager() ServerManager {
	return ServerManager{
		Servers: make(map[string]Server),
	}
}

func NewServer(name, kind, username, password string) (Server, error) {
	var s Server
	var err error

	switch kind {
	case "megarac":
		s, err = megarac.NewApi(name, true)
		if err != nil {
			return nil, err
		}
		break
	case "hpilo4":
		s, err = hpilo4.NewApi(name, true)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown server kind %s", kind)
	}
	err = s.Login(username, password)
	if err != nil {
		return nil, err
	}
	return s, nil
}
