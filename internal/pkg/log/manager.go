package log

import (
	"strings"

	"github.com/kristofferahl/go-centry/internal/pkg/io"
	"github.com/sirupsen/logrus"
)

// Config encapsulates configuration for the logger factory
type Config struct {
	Level  string
	Prefix string
	IO     io.InputOutput
}

// Manager handles creations of loggers
type Manager struct {
	config *Config
	logger *logrus.Logger
}

// CreateManager creates a new logger factory
func CreateManager(level string, prefix string, io io.InputOutput) *Manager {
	return &Manager{
		config: &Config{
			Level:  level,
			Prefix: prefix,
			IO:     io,
		},
		logger: nil,
	}
}

// GetLogger creates a new logger
func (m *Manager) GetLogger() *logrus.Logger {
	if m.logger == nil {
		m.logger = logrus.New()
		m.logger.Out = m.config.IO.Stderr // TODO: Allow changing to Stdout

		l, _ := logrus.ParseLevel(m.config.Level)
		m.logger.SetLevel(l)

		m.logger.Formatter = &PrefixedTextFormatter{
			Prefix: m.config.Prefix,
		}
	}

	return m.logger
}

// TrySetLogLevel tries to change the log level for new loggers
func (m *Manager) TrySetLogLevel(level string) {
	if !strings.EqualFold(level, m.config.Level) {
		logger := m.GetLogger()
		logger.Debugf("changing loglevel to %s (from %s)", level, m.config.Level)
		l, _ := logrus.ParseLevel(level)
		logger.SetLevel(l)
		logger.Debugf("changed loglevel to %s (from %s)", l, m.config.Level)
		m.config.Level = l.String()
	}
}
