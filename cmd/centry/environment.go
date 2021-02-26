package main

import (
	"fmt"
	"os"
	"strings"
)

func environmentOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func overrideFromEnvironment(runtime *Runtime) {
	context := runtime.context
	cli := runtime.cli
	envVersionName := fmt.Sprintf("%s_VERSION", strings.ToUpper(context.manifest.Config.Name))
	envVersion := os.Getenv(envVersionName)
	if envVersion != "" {
		runtime.events = append(runtime.events, fmt.Sprintf("setting the version from environment variable \"%s\"", envVersionName))
		cli.Version = envVersion
	}
}
