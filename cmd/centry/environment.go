package main

import (
	"fmt"
	"os"
	"strings"
)

func overrideFromEnvironment(runtime *Runtime) {
	context := runtime.context
	cli := runtime.cli
	envVersionName := fmt.Sprintf("%s_VERSION", strings.ToUpper(context.manifest.Config.Name))
	envVersion := os.Getenv(envVersionName)
	if envVersion != "" {
		runtime.events = append(runtime.events, fmt.Sprintf("Setting the version from environment variable \"%s\"", envVersionName))
		cli.Version = envVersion
	}
}
