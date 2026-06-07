package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	KeycloakURL       string `yaml:"keycloak_url"`
	ClientID          string `yaml:"client_id"`
	ClientSecret      string `yaml:"client_secret"`
	Password          string `yaml:"password"`
	Username          string `yaml:"username"`
	MatrixURL         string `yaml:"matrix_url"`
	MacaroonSercetKey string `yaml:"macaroon_secret_key"`
}

func initConfiguration() Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return config
}
