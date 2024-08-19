package config

import (
	"fmt"
	"os"

	"github.com/andrewjjenkins/powerlab/pkg/server"

	"gopkg.in/yaml.v3"
)

type ServerCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Kind     string `"yaml:"kind"`
}

type Credentials struct {
	Servers map[string]ServerCredentials `yaml:"credentials"`
}

func LoadCredentials(filename string) (*Credentials, error) {
	if filename == "" {
		filename = "credentials.yaml"
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed open credentials %s: %v", filename, err)
	}

	c := Credentials{}
	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("failed parsing credentials %s: %v", filename, err)
	}

	// Set default server kind
	for k, creds := range c.Servers {
		if creds.Kind == "" {
			creds.Kind = "megarac"
		}
		c.Servers[k] = creds
	}

	return &c, nil
}

func LoadServers(credentialsFilename string) (*server.ServerManager, error) {
	credentials, err := LoadCredentials(credentialsFilename)
	if err != nil {
		return nil, err
	}

	manager := server.NewServerManager()

	for k, c := range credentials.Servers {
		server, err := server.NewServer(k, c.Kind, c.Username, c.Password)
		if err != nil {
			return nil, err
		}
		manager.Servers[k] = server
	}

	return &manager, nil
}
