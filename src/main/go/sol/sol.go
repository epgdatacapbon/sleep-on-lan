package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
)

var logger Logger
var configuration = Configuration{}
var configurationFileName = "sol.json"

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	for _, listenerConfiguration := range configuration.listenersConfiguration {
		if listenerConfiguration.active {
			if strings.EqualFold(listenerConfiguration.nature, "UDP") {
				go ListenerUDP(listenerConfiguration.port)
			} else if strings.EqualFold(listenerConfiguration.nature, "HTTP") {
				go ListenerHTTP(listenerConfiguration.port)
			}
		}
	}
}
func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	fmt.Println(application.Name, application.Version)

	svcConfig := &service.Config{
		Name:        "SleepOnLan",
		DisplayName: application.Name,
		Description: "This service allows a PC to be put into sleep.",
	}

	prg := &program{}
	srv, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		err = service.Control(srv, os.Args[1])
		if err != nil {
			fmt.Println("Failed ("+os.Args[1]+"): ", err)
		} else {
			fmt.Println("Succeeded (" + os.Args[1] + ")")
		}
		os.Exit(0)
	}

	logger.logLevel = 1
	logger.srvLogger, err = srv.Logger(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configuration.InitDefaultConfiguration()
	configuration.Load(dir + string(os.PathSeparator) + configurationFileName)
	configuration.Parse()

	err = srv.Run()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
