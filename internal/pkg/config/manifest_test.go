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
			_, err := LoadManifest("test/data/manifest_test_invalid.yaml")
			g.Assert(err != nil).IsTrue("expected validation error")
		})

		g.It("returns manifest when file is found", func() {
			path := "test/data/manifest_test_valid.yaml"
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

	g.Describe("read file", func() {
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

	g.Describe("parse file", func() {
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

	g.Describe("command", func() {
		g.Describe("annotations", func() {
			g.It("returns nil when command has no annotations", func() {
				c := Command{}
				annotation, err := c.Annotation("x", "y")
				g.Assert(annotation == nil).IsTrue("exected no annotation")
				g.Assert(err == nil).IsTrue("expected error to be nil")
			})

			g.It("returns nil when command does not contain annotation", func() {
				c := Command{
					Annotations: map[string]string{
						"centry.foo/bar": "baz",
					},
				}
				annotation, err := c.Annotation("x", "y")
				g.Assert(annotation == nil).IsTrue("exected no annotation")
				g.Assert(err == nil).IsTrue("expected error to be nil")
			})

			g.It("returns annotation when command contains annotation", func() {
				c := Command{
					Annotations: map[string]string{
						"centry.x/x": "1",
						"centry.y/y": "2",
						"centry.z/z": "3",
					},
				}
				annotation, err := c.Annotation("centry.y", "y")
				g.Assert(annotation != nil).IsTrue("exected annotation")
				g.Assert(annotation.Namespace).Equal("centry.y")
				g.Assert(annotation.Key).Equal("y")
				g.Assert(annotation.Value).Equal("2")
				g.Assert(err == nil).IsTrue("expected error to be nil")
			})
		})
	})

	g.Describe("option", func() {
		g.Describe("annotations", func() {
			g.It("returns nil when option has no annotations", func() {
				o := Option{}
				annotation, err := o.Annotation("x", "y")
				g.Assert(annotation == nil).IsTrue("exected no annotation")
				g.Assert(err == nil).IsTrue("expected error to be nil")
			})

			g.It("returns nil when option does not contain annotation", func() {
				o := Option{
					Annotations: map[string]string{
						"centry.foo/bar": "baz",
					},
				}
				annotation, err := o.Annotation("x", "y")
				g.Assert(annotation == nil).IsTrue("exected no annotation")
				g.Assert(err == nil).IsTrue("expected error to be nil")
			})

			g.It("returns annotation when option contains annotation", func() {
				o := Option{
					Annotations: map[string]string{
						"centry.x/x": "1",
						"centry.y/y": "2",
						"centry.z/z": "3",
					},
				}
				annotation, err := o.Annotation("centry.y", "y")
				g.Assert(annotation != nil).IsTrue("exected annotation")
				g.Assert(annotation.Namespace).Equal("centry.y")
				g.Assert(annotation.Key).Equal("y")
				g.Assert(annotation.Value).Equal("2")
				g.Assert(err == nil).IsTrue("expected error to be nil")
			})
		})
	})
}
