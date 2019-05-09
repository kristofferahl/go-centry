package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type manifest struct {
	Scripts  []string  `yaml:"scripts,omitempty"`
	Commands []command `yaml:"commands,omitempty"`
	Options  []option  `yaml:"options,omitempty"`
	Config   config    `yaml:"config,omitempty"`
	Path     string
	BasePath string
}

type command struct {
	Name        string `yaml:"name,omitempty"`
	Path        string `yaml:"path,omitempty"`
	Help        string `yaml:"help,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type option struct {
	Name        string `yaml:"name,omitempty"`
	EnvName     string `yaml:"env_name,omitempty"`
	Default     string `yaml:"default,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type config struct {
	Name    string    `yaml:"name,omitempty"`
	Version string    `yaml:"version,omitempty"`
	Log     logConfig `yaml:"log,omitempty"`
}

type logConfig struct {
	Level  string `yaml:"level,omitempty"`
	Prefix string `yaml:"prefix,omitempty"`
}

func loadManifest(path string) *manifest {
	mp, _ := filepath.Abs(path)

	if _, err := os.Stat(mp); os.IsNotExist(err) {
		fmt.Println("The first argument of centry must be a path to a valid manfest file")
		os.Exit(1)
	}

	bs := readManifest(mp)
	m := parseManifest(bs)

	m.Path = mp
	m.BasePath = filepath.Dir(mp)
	return m
}

func readManifest(filename string) []byte {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to load manifest file.", "Error:", err)
		os.Exit(1)
	}
	return bs
}

func parseManifest(bs []byte) *manifest {
	m := manifest{}
	err := yaml.Unmarshal(bs, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &m
}
