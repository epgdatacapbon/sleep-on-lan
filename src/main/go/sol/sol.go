package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
)

var configuration = Configuration{}
var configurationFileName = "sol.json"

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	for _, listenerConfiguration := range configuration.listenersConfiguration {
		if listenerConfiguration.active {
			if strings.EqualFold(listenerConfiguration.nature, "UDP") {
				go ListenUDP(listenerConfiguration.port)
			} else if strings.EqualFold(listenerConfiguration.nature, "HTTP") {
				go ListenHTTP(listenerConfiguration.port) // , configuration.Commands, configuration.Auth, configuration.HTTPOutput)
			}
		}
	}
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	PreInitLoggers()
	Info.Println(Version.ApplicationName + " Version " + Version.Version())

	svcConfig := &service.Config{
		Name:        "SleepOnLan",
		DisplayName: Version.ApplicationName,
		Description: "This service allows a PC to be put into sleep.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		Error.Println(err)
		return;
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		Error.Println(err)
		return;
	}

	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			Error.Println("Failed (" + os.Args[1] + ") :", err)
		} else {
			Info.Println("Succeeded (" + os.Args[1] + ")")
		}
		return
	}

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configuration.InitDefaultConfiguration()
	configuration.Load(dir + string(os.PathSeparator) + configurationFileName)
	configuration.Parse()

	Info.Println("Hardware IP/mac addresses are : ")
	for key, value := range LocalNetworkMap() {
		Info.Println(" - local IP adress [" + key + "], mac [" + value + "]")
	}

	err = s.Run()
	if err != nil {
		Error.Println(err)
	}
}
