package main

import (
	"os/exec"
	"strings"
	"syscall"

	winio "github.com/Microsoft/go-winio"
)

const (
	DEFAULT_COMMAND_SLEEP     = "sleep"
	DEFAULT_COMMAND_HIBERNATE = "hibernate"
	DEFAULT_COMMAND_SHUTDOWN  = "shutdown"
)

func RegisterDefaultCommand() {
	defaultSleepCommand := CommandConfiguration{Operation: DEFAULT_COMMAND_SLEEP, CommandType: COMMAND_TYPE_INTERNAL}
	defaultHibernateCommand := CommandConfiguration{Operation: DEFAULT_COMMAND_HIBERNATE, CommandType: COMMAND_TYPE_INTERNAL}
	defaultShutdownCommand := CommandConfiguration{Operation: DEFAULT_COMMAND_SHUTDOWN, CommandType: COMMAND_TYPE_INTERNAL}
	configuration.Commands = []CommandConfiguration{defaultSleepCommand, defaultHibernateCommand, defaultShutdownCommand}
}

func ExecuteCommand(Command CommandConfiguration) {
	if Command.CommandType == COMMAND_TYPE_INTERNAL {
		logger(3, "Executing operation ["+Command.Operation+"], type ["+Command.CommandType+"]")
		if Command.Operation == DEFAULT_COMMAND_SLEEP {
			sleepDLLImplementation(0)
		} else if Command.Operation == DEFAULT_COMMAND_HIBERNATE {
			sleepDLLImplementation(1)
		} else if Command.Operation == DEFAULT_COMMAND_SHUTDOWN {
			shutdownDLLImplementation()
		}
	} else if Command.CommandType == COMMAND_TYPE_EXTERNAL {
		logger(3, "Execute operation ["+Command.Operation+"], type ["+Command.CommandType+"], command ["+Command.Command+"]")
		commandImplementation(Command.Command)
	} else {
		logger(2, "Invalid command type ["+Command.CommandType+"]")
	}
}

func sleepDLLImplementation(state int) {
	var mod = syscall.NewLazyDLL("Powrprof.dll")
	var proc = mod.NewProc("SetSuspendState")

	// DLL API : public static extern bool SetSuspendState(bool hiberate, bool forceCritical, bool disableWakeEvent);
	// ex. : uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Done Title"))),
	r, _, _ := proc.Call(
		uintptr(state), // hibernate
		uintptr(0),     // forceCritical
		uintptr(0))     // disableWakeEvent
	if r == 0 {
		logger(2, "Unable to execute Suspend command")
	}
}

func shutdownDLLImplementation() {
	// SeShutdownPrivilege
	err := winio.RunWithPrivilege("SeShutdownPrivilege", func() error {
		var mod = syscall.NewLazyDLL("Advapi32.dll")
		var proc = mod.NewProc("InitiateSystemShutdownW")

		// DLL API : public static extern bool InitiateSystemShutdown(string lpMachineName, string lpMessage, int dwTimeout, bool bForceAppsClosed, bool bRebootAfterShutdown);
		// ex. : uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Done Title"))),

		proc.Call(
			uintptr(0), // lpMachineName
			uintptr(0), // lpMessage
			uintptr(0), // dwTimeout
			uintptr(1), // bForceAppsClosed
			uintptr(0)) // bRebootAfterShutdown
		return nil
	})
	if err != nil {
		logger(2, "Unable to execute shutdown command")
	}
}

func commandImplementation(command string) {
	if command == "" {
		return
	}

	parts := strings.Fields(command)
	head := parts[0]
	parts = parts[1:len(parts)]
	_, err := exec.Command(head, parts...).Output()
	if err != nil {
		logger(2, "Unable to execute command ["+command+"]: "+err.Error())
	}
}
