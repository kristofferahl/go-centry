package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/franela/goblin"
)

func TestManifest(t *testing.T) {
	g := Goblin(t)

	g.Describe("LoadManifest", func() {
		g.It("returns error for invalid manifest file", func() {
			_, err := LoadManifest("test/data/invalid.yaml")
			g.Assert(err != nil).IsTrue("expected validation error")
		})

		g.It("returns manifest when file is found", func() {
			path := "test/data/main_test.yaml"
			absPath, _ := filepath.Abs(path)
			basePath := filepath.Dir(absPath)

			m, err := LoadManifest(path)

			g.Assert(m != nil).IsTrue("exected manifest")
			g.Assert(m.Path).Equal(absPath)
			g.Assert(m.BasePath).Equal(basePath)
			g.Assert(err == nil).IsTrue("expected error to be nil")
		})

		g.It("returns error when path is invalid", func() {
			m, err := LoadManifest("foo")
			g.Assert(m == nil).IsTrue("exected manifest to be nil")
			g.Assert(err != nil).IsTrue("expected error")
			g.Assert(err.Error()).Equal("The first argument must be a path to a valid manifest file (foo)")
		})
	})

	g.Describe("readManifestFile", func() {
		g.It("returns byte slice when file is found", func() {
			file, _ := ioutil.TempFile("", "manifest")
			defer os.Remove(file.Name())
			bs, err := readManifestFile(file.Name())
			g.Assert(bs != nil).IsTrue("exected byte slice")
			g.Assert(err == nil).IsTrue("expected error to be nil")
		})

		g.It("returns error when file not found", func() {
			bs, err := readManifestFile("foo")
			g.Assert(bs == nil).IsTrue("exected byte slice to be nil")
			g.Assert(err != nil).IsTrue("expected error")
			g.Assert(strings.HasPrefix(err.Error(), "Failed to read manifest file")).IsTrue("expected error message")
		})
	})

	g.Describe("parseManifestFile", func() {
		g.It("returns manifest for valid yaml", func() {
			m, err := parseManifestYaml([]byte(`config:`))
			g.Assert(m != nil).IsTrue("exected manifest")
			g.Assert(err == nil).IsTrue("expected error to be nil")
		})

		g.It("returns error when byte slice in invalid yaml", func() {
			m, err := parseManifestYaml([]byte("invalid yaml"))
			g.Assert(m == nil).IsTrue("exected manifest to be nil")
			g.Assert(err != nil).IsTrue("expected error")
			g.Assert(strings.HasPrefix(err.Error(), "Failed to parse manifest yaml")).IsTrue("expected error message")
		})
	})
}
