package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type config struct {
	Title           string
	Description     string
	Keywords		string
	Posts           string
	Templates       string
	Public          string
	PreviewLength   int
	PostsPerPage    int
	DisqusShortname string
	GoogleAnalyticsID string
}

func loadConfig(configfile string) config {
	if _, err := os.Stat(configfile); err != nil {
		log.Fatal("Config file not found: ", configfile)
	}

	var cfg config
	if _, err := toml.DecodeFile(configfile, &cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
