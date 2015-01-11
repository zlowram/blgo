package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	Title           string
	Description     string
	Posts           string
	Templates       string
	Public          string
	PreviewLength   int
	PostsPerPage    int
	DisqusShortname string
}

func loadConfig(configfile string) Config {
	// Check if config file is present
	if _, err := os.Stat(configfile); err != nil {
		log.Fatal("Config file not found: ", configfile)
	}

	// Load it!
	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
