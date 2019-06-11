package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var s string

	// Basic authentication
	if !configuration.Auth.isEmpty() {
		username, password, ok := r.BasicAuth()
		if ok == false || username != configuration.Auth.Login || password != configuration.Auth.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="auth area"`)
			http.Error(w, "Authorization failed", http.StatusUnauthorized)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, application.Name, application.Version)

	operation := r.URL.Path[1:]
	if operation == "" {
		fmt.Fprintln(w, "\nHosts:")
		interfaces := LocalNetworkMap()
		ips := make([]string, 0, len(interfaces))
		for key := range interfaces {
			ips = append(ips, key)
		}
		sort.Strings(ips)
		for _, ip := range ips {
			fmt.Fprintf(w, `  ip="%s" mac="%s"`+"\n", ip, interfaces[ip])
		}
		fmt.Fprintln(w, "\nListeners:")
		for _, listenerConfiguration := range configuration.listenersConfiguration {
			fmt.Fprintf(w, `  type="%s" port="%d" active="%t"`+"\n", listenerConfiguration.nature,
				listenerConfiguration.port, listenerConfiguration.active)
		}
		fmt.Fprintln(w, "\nCommands:")
		for _, commandConfiguration := range configuration.Commands {
			fmt.Fprintf(w, `  operation="%s" command="%s" type="%s"`+"\n", commandConfiguration.Operation,
				commandConfiguration.Command, commandConfiguration.CommandType)
		}
	} else if len(operation) > 3 && operation[:4] == "wol/" {
		mac := operation[4:]
		if mac == "" {
			s = "No MAC address"
			fmt.Fprintln(w, s)
			logger.Warning(s)
			return
		}
		magicPacket, err := EncodeMagicPacket(mac)
		if err != nil {
			fmt.Fprintln(w, err)
			logger.Warning(err.Error())
			return
		}
		s = "Sending a magic packet to MAC address [" + mac + "]"
		fmt.Fprintln(w, s)
		logger.Info(s)
		err = magicPacket.Wake(configuration.BroadcastIP)
		if err == nil {
		} else {
			fmt.Fprintln(w, err)
			logger.Warning(err.Error())
		}
	} else if operation == "quit" {
		logger.Info("Quit")
		defer os.Exit(0)
	} else {
		for idx, _ := range configuration.Commands {
			availableCommand := configuration.Commands[idx]
			if availableCommand.Operation == operation {
				defer ExecuteCommand(availableCommand)
				return
			}
		}
		s = "Invalid operation [" + operation + "]"
		fmt.Fprintln(w, s)
		logger.Warning(s)
	}
}

func ListenerHTTP(port int) {
	logger.Info("Listening HTTP requests on port [" + strconv.Itoa(port) + "]")
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		logger.Error("Error while start listening: " + err.Error())
		os.Exit(1)
	}
}
