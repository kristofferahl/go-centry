package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/kristofferahl/go-centry/internal/pkg/cmd"
	yaml2 "gopkg.in/yaml.v2"
)

// Manifest defines the structure of a manifest
type Manifest struct {
	Scripts  []string  `yaml:"scripts,omitempty"`
	Commands []Command `yaml:"commands,omitempty"`
	Options  []Option  `yaml:"options,omitempty"`
	Config   Config    `yaml:"config,omitempty"`
	Path     string
	BasePath string
}

// Command defines the structure of commands
type Command struct {
	Name        string            `yaml:"name,omitempty"`
	Path        string            `yaml:"path,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Help        string            `yaml:"help,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Hidden      bool              `yaml:"hidden,omitempty"`
}

// Annotation returns a parsed annotation if present
func (c Command) Annotation(namespace, key string) (*Annotation, error) {
	return ParseAnnotation(getAnnotationString(c.Annotations, namespace, key))
}

// Option defines the structure of options
type Option struct {
	Type        cmd.OptionType    `yaml:"type,omitempty"`
	Name        string            `yaml:"name,omitempty"`
	Short       string            `yaml:"short,omitempty"`
	EnvName     string            `yaml:"env_name,omitempty"`
	Default     string            `yaml:"default,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Hidden      bool              `yaml:"hidden,omitempty"`
}

// Annotation returns a parsed annotation if present
func (o Option) Annotation(namespace, key string) (*Annotation, error) {
	return ParseAnnotation(getAnnotationString(o.Annotations, namespace, key))
}

// Config defines the structure for the configuration section
type Config struct {
	Name                 string    `yaml:"name,omitempty"`
	Description          string    `yaml:"description,omitempty"`
	Version              string    `yaml:"version,omitempty"`
	Log                  LogConfig `yaml:"log,omitempty"`
	EnvironmentPrefix    string    `yaml:"environmentPrefix,omitempty"`
	HideInternalCommands bool      `yaml:"hideInternalCommands,omitempty"`
	HideInternalOptions  bool      `yaml:"hideInternalOptions,omitempty"`
}

// LogConfig defines the structure for log configuration section
type LogConfig struct {
	Level  string `yaml:"level,omitempty"`
	Prefix string `yaml:"prefix,omitempty"`
}

// LoadManifest reads, parses and returns a manifest root object
func LoadManifest(manifest string) (*Manifest, error) {
	mp, _ := filepath.Abs(manifest)

	if _, err := os.Stat(mp); os.IsNotExist(err) {
		return nil, fmt.Errorf("manifest file not found (path=%s)", manifest)
	}

	bs, err := readManifestFile(mp)
	if err != nil {
		return nil, err
	}

	m, err := parseManifestYaml(bs)
	if err != nil {
		return nil, err
	}

	jbs, err := yaml.YAMLToJSON(bs)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(jbs)
	err = validateManifestYaml("bindata://schemas/manifest.json", r)
	if err != nil {
		return nil, err
	}

	m.Path = mp
	m.BasePath = filepath.Dir(mp)

	return m, nil
}

func readManifestFile(filename string) ([]byte, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file (path=%s). %v", filename, err)
	}
	return bs, nil
}

func parseManifestYaml(bs []byte) (*Manifest, error) {
	m := Manifest{}
	err := yaml2.Unmarshal(bs, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manifest yaml. %v", err)
	}
	return &m, nil
}

func getAnnotationString(annotations map[string]string, namespace, key string) string {
	if annotations == nil {
		return ""
	}

	namespaceKey := AnnotationNamespaceKey(namespace, key)
	value := annotations[namespaceKey]
	if value == "" {
		return ""
	}

	return AnnotationString(namespace, key, value)
}
