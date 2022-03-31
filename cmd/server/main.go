package main

import (
	"flag"
	"log"
	"stochastic_indicator/internal"
	"stochastic_indicator/internal/config"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/service.toml", "path to config")
}

func main() {
	flag.Parse()

	configuration := config.NewConfig()
	if _, err := toml.DecodeFile(configPath, configuration); err != nil {
		log.Fatalln(err)
	}

	if err := internal.Start(configuration); err != nil {
		log.Fatalln(err)
	}
}
