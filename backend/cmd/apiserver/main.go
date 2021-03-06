package main

import (
	"flag"
	"log"

	"github.com/MrSedan/restapigoown/backend/internal/app/apiserver"
	"github.com/pelletier/go-toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "/go/bin/configs/server.toml", "path to config")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	data, err := toml.LoadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	data.Unmarshal(config)
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
