package main

func logger(logType int, s string) {
	if logType <= logLevel {
		switch logType {
		case 1:
			srvLogger.Error(s)
		case 2:
			srvLogger.Warning(s)
		case 3:
			srvLogger.Info(s)
		default:
			return
		}
	}
}
