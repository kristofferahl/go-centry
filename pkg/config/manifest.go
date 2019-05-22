package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
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
	Help        string            `yaml:"help,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// Option defines the structure of options
type Option struct {
	Name        string            `yaml:"name,omitempty"`
	EnvName     string            `yaml:"env_name,omitempty"`
	Default     string            `yaml:"default,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// Config defines the structure for the configuration section
type Config struct {
	Name    string    `yaml:"name,omitempty"`
	Version string    `yaml:"version,omitempty"`
	Log     LogConfig `yaml:"log,omitempty"`
}

// LogConfig defines the structure for log configuration section
type LogConfig struct {
	Level  string `yaml:"level,omitempty"`
	Prefix string `yaml:"prefix,omitempty"`
}

// LoadManifest reads, parses and returns a manifest root object
func LoadManifest(path string) *Manifest {
	mp, _ := filepath.Abs(path)

	if _, err := os.Stat(mp); os.IsNotExist(err) {
		fmt.Println("The first argument must be a path to a valid manfest file")
		os.Exit(1)
	}

	bs := readManifestFile(mp)
	m := parseYaml(bs)

	m.Path = mp
	m.BasePath = filepath.Dir(mp)
	return m
}

func readManifestFile(filename string) []byte {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to load manifest file.", "Error:", err)
		os.Exit(1)
	}
	return bs
}

func parseYaml(bs []byte) *Manifest {
	m := Manifest{}
	err := yaml.Unmarshal(bs, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &m
}
