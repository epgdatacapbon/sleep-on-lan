package main

import (
	"encoding/xml"
	"net/http"
	"sort"
	"strconv"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type RestResultHost struct {
	XMLName    xml.Name `xml:"host" json:"-"`
	Ip         string   `xml:"ip,attr"`
	MacAddress string   `xml:"mac,attr"`
}

type RestResultHosts struct {
	XMLName xml.Name `xml:"hosts" json:"-"`
	Hosts   []RestResultHost
}

type RestResultCommands struct {
	XMLName  xml.Name `xml:"commands" json:"-"`
	Commands []RestResultCommandConfiguration
}

type RestResultCommandConfiguration struct {
	XMLName   xml.Name `xml:"command" json:"-"`
	Operation string   `xml:"operation,attr"`
	Command   string   `xml:"command,attr"`
	Type      string   `xml:"type,attr"`
}

type RestResultListeners struct {
	XMLName   xml.Name `xml:"listeners" json:"-"`
	Listeners []RestResultListenerConfiguration
}

type RestResultListenerConfiguration struct {
	XMLName xml.Name `xml:"listener" json:"-"`
	Type    string   `xml:"type,attr"`
	Port    int      `xml:"port,attr"`
	Active  bool     `xml:"active,attr"`
}

type RestResult struct {
	XMLName     xml.Name `xml:"result" json:"-"`
	Application string   `xml:"application"`
	Version     string   `xml:"version"`
	Hosts       RestResultHosts
	Listeners   RestResultListeners
	Commands    RestResultCommands
}

type RestOperationResult struct {
	XMLName   xml.Name `xml:"result" json:"-"`
	Operation string   `xml:"operation"`
	Result    bool     `xml:"successful"`
}

func renderResult(c echo.Context, status int, result interface{}) error {
	format := c.QueryParam("format")
	if strings.EqualFold(configuration.HTTPOutput, "JSON") || strings.EqualFold(format, "JSON") {
		return c.JSONPretty(status, result, "  ")
	} else {
		return c.XMLPretty(status, result, "  ")
	}
}

func ListenHTTP(port int) {
	e := echo.New()
	e.HideBanner = true

	if !configuration.Auth.isEmpty() {
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			if username == configuration.Auth.Login && password == configuration.Auth.Password {
				return true, nil
			}
			return false, nil
		}))
	}

	e.GET("/", func(c echo.Context) error {
		result := &RestResult{}
		result.Application = Version.ApplicationName
		result.Version = Version.Version()
		result.Hosts = RestResultHosts{}
		result.Listeners = RestResultListeners{}
		result.Commands = RestResultCommands{}

		interfaces := LocalNetworkMap()
		ips := make([]string, 0, len(interfaces))
		for key := range interfaces {
			ips = append(ips, key)
		}
		sort.Strings(ips)
		for _, ip := range ips {
			result.Hosts.Hosts = append(result.Hosts.Hosts, RestResultHost{Ip: ip, MacAddress: interfaces[ip]})
		}
		for _, listenerConfiguration := range configuration.listenersConfiguration {
			result.Listeners.Listeners = append(result.Listeners.Listeners, RestResultListenerConfiguration{Type: listenerConfiguration.nature, Port: listenerConfiguration.port, Active: listenerConfiguration.active})
		}

		for _, commandConfiguration := range configuration.Commands {
			result.Commands.Commands = append(result.Commands.Commands, RestResultCommandConfiguration{Type: commandConfiguration.CommandType, Operation: commandConfiguration.Operation, Command: commandConfiguration.Command})
		}

		return renderResult(c, http.StatusOK, result)
	})

	// N.B.: sleep operation is now registred through commands below
	for _, command := range configuration.Commands {
		e.GET("/" + command.Operation, func(c echo.Context) error {
			
			items := strings.Split(c.Request().URL.Path, "/")
			operation := items[1]

			result := &RestOperationResult{
				Operation:  operation,
				Result: true,
			}
			for idx, _ := range configuration.Commands {
				availableCommand := configuration.Commands[idx]
				if availableCommand.Operation == operation {
					logger.Info("Executing [" + operation + "]")
					ExecuteCommand(availableCommand)
					break
				}
			}
			return renderResult(c, http.StatusOK, result)
		})
	}

	e.GET("/wol/:mac", func(c echo.Context) error {
		result := &RestOperationResult{
			Operation:  "wol",
			Result: true,
		}

		mac := c.Param("mac")
		logger.Info("Sending wol magic packet to MAC address [" + mac + "]")
		magicPacket, err := EncodeMagicPacket(mac)
		if err != nil {
			logger.Error(err)
		} else {
			magicPacket.Wake(configuration.BroadcastIP)
		}
		return renderResult(c, http.StatusOK, result)
	})

	e.GET("/quit", func(c echo.Context) error {
		result := &RestOperationResult{
			Operation:  "quit",
			Result: true,
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextXMLCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)		
		b, _ := xml.Marshal(result)
		c.Response().Write(b)
		c.Response().Flush()
		defer os.Exit(1)
		return nil
	})

	err := e.Start(":" + strconv.Itoa(port))
	if err != nil {
		logger.Error("Error while start listening: ", err)
	}
}
