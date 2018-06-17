package main

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"

	"github.com/sniperkit/cxdig/pkg/cmd"
)

// HandleCleanExit makes sure to exit properly the application
func HandleCleanExit() {
	if r := recover(); r != nil {
		stackTrace := string(debug.Stack())
		logrus.WithField("stack", stackTrace).Fatalf("PANIC: %v", r)
	}
}

func main() {
	defer HandleCleanExit()

	cmd.Execute()
}
