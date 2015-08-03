package main

import (
  "flag"
)

type Settings struct {
	logLevel *int
	logPath *string
	port *string
	domain *string
	secure *bool
  auth *string
}

var settings *Settings

func parseFlags() {
  settings = &Settings{}
  settings.logLevel = flag.Int("log-level", ERROR, "log level 0 - 15")
  settings.logPath = flag.String("log-path", "", "path to file where logs are stored")
  settings.port = flag.String("port", ":8080", "server port")
  settings.domain = flag.String("domain", "localhost", "server domain")
  settings.secure = flag.Bool("secure", false, "use secure urls (wss, https protocols)")
  settings.auth = flag.String("auth", "", "provide basic auth username and password")
  flag.Parse()
}
