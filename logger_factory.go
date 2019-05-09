package main

import (
	"strings"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type loggerConfig struct {
	level  string
	prefix string
}

type loggerFactory struct {
	config *loggerConfig
}

func (f *loggerFactory) createLogger() *logrus.Logger {
	c := f.config
	logger = logrus.New()

	l, _ := logrus.ParseLevel(c.level)
	logger.SetLevel(l)

	logger.Formatter = &PrefixedTextFormatter{
		Prefix: c.prefix,
	}
	return logger
}

func (f *loggerFactory) trySetLogLevel(level string) {
	if !strings.EqualFold(level, f.config.level) {
		logger.Debugf("Changing loglevel to %s (from %s)", level, f.config.level)
		l, _ := logrus.ParseLevel(level)
		logger.SetLevel(l)
		logger.Debugf("Changed loglevel to %s (from %s)", l, f.config.level)
		f.config.level = l.String()
	}
}
