package main

import (
  "flag"
)

type Settings struct {
	logLevel *int
	logPath *string
}

var settings *Settings

func parseFlags() {
  settings = &Settings{}
  settings.logLevel = flag.Int("log-level", ERROR, "log level 0 - 15")
  settings.logPath = flag.String("log-path", "", "path to file where logs are stored")
  flag.Parse()
}
