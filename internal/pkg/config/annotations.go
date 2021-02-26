package config

import (
	"fmt"
	"regexp"
	"strings"
)

// TrueString defines a true value
const TrueString string = "true"

// CommandAnnotationCmdNamespace defines an annotation namespace
const CommandAnnotationCmdNamespace string = "centry.cmd"

// CommandAnnotationCmdOptionNamespace defines an annotation namespace
const CommandAnnotationCmdOptionNamespace string = "centry.cmd.option"

// CommandAnnotationAPINamespace defines an annotation namespace
const CommandAnnotationAPINamespace string = "centry.api"

// Annotation defines an annotation
type Annotation struct {
	Namespace       string
	NamespaceValues map[string]string
	Key             string
	Value           string
}

// AnnotationNamespaceKey creates an annotation namespace/key string
func AnnotationNamespaceKey(namespace, key string) string {
	return fmt.Sprintf("%s/%s", namespace, key)
}

// AnnotationString creates an annotation string
func AnnotationString(namespace, key, value string) string {
	return fmt.Sprintf("%s=%s", AnnotationNamespaceKey(namespace, key), value)
}

// ParseAnnotation parses text into an annotation
func ParseAnnotation(text string) (*Annotation, error) {
	text = strings.TrimSpace(text)

	if !strings.HasPrefix(text, "centry") || !strings.Contains(text, "/") || !strings.Contains(text, "=") {
		return nil, nil
	}

	namespaceKeyValueString := strings.SplitN(text, "/", 2)
	namespace := namespaceKeyValueString[0]
	keyValueString := namespaceKeyValueString[1]

	kvp := strings.SplitN(keyValueString, "=", 2)
	if len(kvp) != 2 {
		return nil, fmt.Errorf("failed to parse annotation! The text \"%s\" is not a valid annotation", text)
	}

	return &Annotation{
		Namespace:       cleanupNamespace(namespace),
		NamespaceValues: extractNamespaceValues(namespace),
		Key:             kvp[0],
		Value:           kvp[1],
	}, nil
}

func extractNamespaceValues(namespace string) (params map[string]string) {
	var compRegEx = regexp.MustCompile("\\.(\\w+)\\[([0-9A-Za-z_:]+)\\]")
	match := compRegEx.FindAllStringSubmatch(namespace, -1)

	params = make(map[string]string)
	for _, kv := range match {
		k := kv[1]
		v := kv[2]
		params[k] = v
	}

	return
}

func cleanupNamespace(namespace string) string {
	var compRegEx = regexp.MustCompile("(\\.*)(\\[([0-9A-Za-z_:]+)\\])")
	return compRegEx.ReplaceAllString(namespace, `${1}`)
}
