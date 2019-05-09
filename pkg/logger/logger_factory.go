package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// Config encapsulates configuration for the logger factory
type Config struct {
	Level  string
	Prefix string
}

// Factory handles creations of loggers
type Factory struct {
	config *Config
}

// CreateFactory creates a new logger factory
func CreateFactory(level string, prefix string) *Factory {
	return &Factory{
		config: &Config{
			Level:  level,
			Prefix: prefix,
		},
	}
}

// CreateLogger creates a new logger
func (f *Factory) CreateLogger() *logrus.Logger {
	c := f.config
	logger = logrus.New()

	l, _ := logrus.ParseLevel(c.Level)
	logger.SetLevel(l)

	logger.Formatter = &PrefixedTextFormatter{
		Prefix: c.Prefix,
	}
	return logger
}

// TrySetLogLevel tries to change the log level for new loggers
func (f *Factory) TrySetLogLevel(level string) {
	if !strings.EqualFold(level, f.config.Level) {
		logger.Debugf("Changing loglevel to %s (from %s)", level, f.config.Level)
		l, _ := logrus.ParseLevel(level)
		logger.SetLevel(l)
		logger.Debugf("Changed loglevel to %s (from %s)", l, f.config.Level)
		f.config.Level = l.String()
	}
}
