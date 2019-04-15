package main

import (
	"log"
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
	var err error
	for _, listenerConfiguration := range configuration.listenersConfiguration {
		if listenerConfiguration.active {
			if strings.EqualFold(listenerConfiguration.nature, "UDP") {
				err = ListenUDP(listenerConfiguration.port)
			} else if strings.EqualFold(listenerConfiguration.nature, "HTTP") {
				err = ListenHTTP(listenerConfiguration.port)
			}
			if err != nil {
				os.Exit(1)
			}
		}
	}
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "SleepOnLan",
		DisplayName: Version.ApplicationName,
		Description: "This service allows a PC to be put into sleep.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Print("Failed (" + os.Args[1] + "): ", err)
		} else {
			log.Print("Succeeded (" + os.Args[1] + ")")
		}
		return
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info(Version.ApplicationName + " Version " + Version.Version())
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configuration.InitDefaultConfiguration()
	configuration.Load(dir + string(os.PathSeparator) + configurationFileName)
	configuration.Parse()

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
