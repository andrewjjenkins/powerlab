package server

import (
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

func NewServer(name, username, password string) (Server, error) {
	// FIXME: Ensure the type is megarac
	s, err := megarac.NewApi(name, true)
	if err != nil {
		return nil, err
	}
	err = s.Login(username, password)
	if err != nil {
		return nil, err
	}
	return s, nil
}
