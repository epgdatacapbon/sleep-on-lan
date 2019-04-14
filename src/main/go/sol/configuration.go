package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

const (
	COMMAND_TYPE_EXTERNAL = "external"
	COMMAND_TYPE_INTERNAL = "internal"
)

type Configuration struct {
	Listeners   []string // what is read from the sol.json configuration file
	BroadcastIP string
	Commands    []CommandConfiguration // the various defined commands. Will be enhanded with default operation if empty from configuration
	Default     string
	Auth        AuthConfiguration  // optional
	HTTPOutput  string

	listenersConfiguration []ListenerConfiguration // converted once parsed from Listeners
}

type AuthConfiguration struct {
	Login 		string `json:"Login"`
	Password	string `json:"password"`
}

func (a AuthConfiguration) isEmpty() bool {
	return a.Login == "" && a.Password == ""
}

type CommandConfiguration struct {
	Operation    string `json:"Operation"`
	Command      string `json:"Command"`
	CommandType  string `json:"Type"`
}

type ListenerConfiguration struct {
	active bool
	port   int
	nature string
}

func (conf *Configuration) InitDefaultConfiguration() {
	conf.Listeners = []string{"UDP:9", "HTTP:8009"}
	conf.BroadcastIP = "192.168.255.255"
	conf.HTTPOutput = "XML"
	// default commands are registered on Parse() method, depending on the current operating system
}

func (conf *Configuration) Load(configurationFileName string) {
	if _, err := os.Stat(configurationFileName); err == nil {
		logger.Info("Configuration file found at [" + configurationFileName + "]")
		file, _ := os.Open(configurationFileName)
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&conf)
		if err != nil {
			logger.Error("error while loading configuration: ", err)
		}
	} else {
		logger.Info("No configuration file found at [" + configurationFileName + "]")
	}
}

func (conf *Configuration) Parse() {
	// Convert activated ports
	for _, s := range conf.Listeners {
		var splitted = strings.Split(s, ":")
		var key = splitted[0]
		var listenerConfiguration = new(ListenerConfiguration)
		listenerConfiguration.active = true
		if len(splitted) == 2 {
			listenerConfiguration.port, _ = strconv.Atoi(splitted[1])
		}
		if strings.EqualFold(key, "UDP") {
			listenerConfiguration.nature = "UDP"
			conf.listenersConfiguration = append(conf.listenersConfiguration, *listenerConfiguration)
		} else if strings.EqualFold(key, "HTTP") {
			listenerConfiguration.nature = "HTTP"
			conf.listenersConfiguration = append(conf.listenersConfiguration, *listenerConfiguration)
		} else {
			logger.Error("Unknown listener type [" + key + "]")
		}
	}

	// If no commands are found, inject default ones
	var nbCommands = len(conf.Commands)
	if nbCommands == 0 {
		RegisterDefaultCommand()
	} else if nbCommands == 1 {
		conf.Default = conf.Commands[0].Operation
	}

	// Set the first command to default if not provided
	if conf.Default == "" {
			conf.Default = conf.Commands[0].Operation
	}
	logger.Info("Set default operation to [" + conf.Default + "]")

	// Set command type
	for idx, _ := range conf.Commands {
		command := &conf.Commands[idx]
		if command.Command == ""  {
			command.CommandType = COMMAND_TYPE_INTERNAL
		} else if command.CommandType == ""  {
			command.CommandType = COMMAND_TYPE_EXTERNAL
		}
	}
}
