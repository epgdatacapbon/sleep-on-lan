package main

func RegisterDefaultCommand() {
	defaultCommand := CommandConfiguration{Operation: "sleep", CommandType: COMMAND_TYPE_EXTERNAL, IsDefault: true, Command: "pm-suspend"}
	configuration.Commands = []CommandConfiguration{defaultCommand}
}

func ExecuteCommand(Command CommandConfiguration) {
	if Command.CommandType == COMMAND_TYPE_EXTERNAL {
		logger.Info("Executing operation [" + Command.Operation + "], type [" + Command.Command + "], command [" + Command.Command + "]")
		sleepCommandLineImplementation(Command.Command)
	} else {
		logger.Error("Unknown command type [" + Command.CommandType + "]")
	}
}

func sleepCommandLineImplementation(cmd string) {
	if cmd == "" {
		cmd = "pm-suspend"
	}
	logger.Info("Sleep implementation [linux], sleep command is [" + cmd + "]")
	_, _, err := Execute(cmd)
	if err != nil {
		logger.Error("Can't execute command [" + cmd + "]: " + err)
	} else {
		logger.Info("Command correctly executed")
	}
}
