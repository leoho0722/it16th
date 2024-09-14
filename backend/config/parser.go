package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"leoho.io/it16th-webauthn-rp-server/utils"
)

var cfg Config

func Parse() Config {
	yamlData, err := os.ReadFile("config/config.yaml")
	if err != nil {
		fmt.Println("cannot read config.yaml")
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		fmt.Println("cannot unmarshal config.yaml")
		panic(err)
	}
	fmt.Println("parse config.yaml success")
	fmt.Println("config: ", utils.PrintJSON(config))

	cfg = config

	return config
}
