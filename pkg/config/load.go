package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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
	return &c, nil
}
