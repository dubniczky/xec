package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var defaultLocalpaths = []string{
	"xec",
	"xec.yaml",
	"xec.json",
	".xec",
	".xec.yaml",
	".xec.json",
}


type Config struct {
	Cmd string
	Args []string
	Stdout string
	Stderr string
}

func readConfig(path string) Config {
	conf := Config{}
	var err error

	configFile, err := os.ReadFile(path)
    if err != nil {
        log.Fatal(err)
    }

	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		log.Fatal(err)
	}

	return conf
}

func parseConfig(data []byte) (Config, error) {
	conf := Config{}
	err := yaml.Unmarshal(data, &conf)
	return conf, err
}

func assembleCommand(conf Config) string {
	cmd := conf.Cmd

	// Arg
	for _, arg := range conf.Args {
		if len(arg) == 0 {
			continue
		}
		if len(arg) == 1 {
			cmd += " -" + arg
			continue
		}
		
		cmd += " --" + arg
	}

	// Stdout & Stderr
	if len(conf.Stdout) > 0 {
		cmd += " > " + conf.Stdout
	}
	if len(conf.Stderr) > 0 {
		cmd += " 2> " + conf.Stderr
	}

	return cmd
}

func main() {
	logerr := log.New(os.Stderr, "ERROR: ", 0)

	// Try default local paths if no argument is provided
	if len(os.Args) < 2 {
		for _, path := range defaultLocalpaths {
			configFile, err := os.ReadFile(path)
			if err == nil {
				log.Println("Found config file at: " + path)
				conf, err := parseConfig(configFile)
				if err != nil {
					logerr.Printf("Could not parse config file: %s", path)
					os.Exit(1)
				}
				fmt.Println(assembleCommand(conf))
				return
			}
		}

		logerr.Println("No default config file found. Please create one or provide a config file path.")
		os.Exit(1)
	}

	// Read config file from argument
	path := os.Args[1]
	conf := readConfig(path)
	fmt.Println(assembleCommand(conf))
	os.Exit(0)
}