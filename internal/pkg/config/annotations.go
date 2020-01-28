package config

import (
	"fmt"
	"strings"
)

// TrueString defines a true value
const TrueString string = "true"

// CommandAnnotationAPIServe defines an annotation
const CommandAnnotationAPIServe string = "centry.api/serve"

// CommandAnnotationDescriptionNamespace denines an annotation namespace
const CommandAnnotationDescriptionNamespace string = "centry.cmd.description"

// CommandAnnotationHelpNamespace denines an annotation namespace
const CommandAnnotationHelpNamespace string = "centry.cmd.help"

// CommandAnnotationNamespaces holds a list of command annotation namespaces
var CommandAnnotationNamespaces = []string{
	CommandAnnotationDescriptionNamespace,
	CommandAnnotationHelpNamespace,
}

// Annotation defines an annotation
type Annotation struct {
	Namespace string
	Key       string
	Value     string
}

// ParseAnnotation parses text into an annotation
func ParseAnnotation(text string) (*Annotation, error) {
	text = strings.TrimSpace(text)

	if !strings.HasPrefix(text, "centry") || !strings.Contains(text, "=") || !strings.Contains(text, "/") {
		return nil, nil
	}

	kvp := strings.Split(text, "=")
	if len(kvp) != 2 {
		return nil, fmt.Errorf("Failed to parse annotation! The text \"%s\" is not a valid annotation", text)
	}

	nk := kvp[0]
	v := kvp[1]

	nkkvp := strings.Split(nk, "/")
	if len(nkkvp) != 2 {
		return nil, fmt.Errorf("Failed to parse annotation! The text \"%s\" is not a valid annotation", text)
	}

	return &Annotation{
		Namespace: nkkvp[0],
		Key:       nkkvp[1],
		Value:     v,
	}, nil
}
