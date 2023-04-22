package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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
	Flags []string
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

	// Flags
	for _, flag := range conf.Flags {
		if len(flag) == 0 {
			continue
		}
		if len(flag) == 1 {
			cmd += " -" + flag
			continue
		}
		
		cmd += " --" + flag
	}

	// Args
	for _, arg := range conf.Args {
		if len(arg) == 0 {
			continue
		}
		
		cmd += " " + arg
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

func start(conf Config, showFlag bool) bool {
	cmd := assembleCommand(conf)
	if showFlag {
		fmt.Println(cmd)
		return false
	}

	cmdWords := strings.Split(cmd, " ")
	command := exec.Command(cmdWords[0], cmdWords[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
    err := command.Run()

    if err != nil {
        log.Fatal(err)
		return false
    }
	return true
}

func main() {
	logerr := log.New(os.Stderr, "ERROR: ", 0)

	var showFlag bool
	flag.BoolVar(&showFlag, "show", false, "Show the command that would be run from the config file")

	flag.Parse()

	// Try default local paths if no argument is provided
	if len(flag.Args()) < 1 {
		for _, path := range defaultLocalpaths {
			configFile, err := os.ReadFile(path)
			if err == nil {
				log.Println("Found config file at: " + path)
				conf, err := parseConfig(configFile)
				if err != nil {
					logerr.Printf("Could not parse config file: %s", path)
					os.Exit(1)
				}
				start(conf, showFlag)
				return
			}
		}

		logerr.Println("No default config file found. Please create one or provide a config file path.")
		os.Exit(1)
	}

	// Read config file from argument
	path := flag.Arg(0)
	conf := readConfig(path)
	start(conf, showFlag)
	os.Exit(0)
}