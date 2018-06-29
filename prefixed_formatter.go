package main

import (
	"github.com/sirupsen/logrus"
)

type PrefixedTextFormatter struct {
	Prefix string
}

func (f *PrefixedTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	formatter := &logrus.TextFormatter{}
	b, err := formatter.Format(entry)
	if err == nil {
		if len(b) > 0 {
			prefix := []byte(f.Prefix)
			buf := append(prefix[:], b[:]...)
			return buf, nil
		}
		return b, nil
	}
	return nil, err
}
