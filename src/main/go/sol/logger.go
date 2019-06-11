package main

import (
	"github.com/kardianos/service"
)

type Logger struct {
	logLevel  int
	srvLogger service.Logger
}

func (l Logger) Error(s string) {
	l.Loggers(1, s)
}

func (l Logger) Warning(s string) {
	l.Loggers(2, s)
}

func (l Logger) Info(s string) {
	l.Loggers(3, s)
}

func (l Logger) Loggers(logType int, s string) {
	if logType <= l.logLevel {
		switch logType {
		case 1:
			l.srvLogger.Error(s)
		case 2:
			l.srvLogger.Warning(s)
		case 3:
			l.srvLogger.Info(s)
		}
	}
}
