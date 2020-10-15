package config

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestAnnotations(t *testing.T) {
	g := Goblin(t)

	g.Describe("ParseAnnotation", func() {
		g.Describe("InvalidAnnotations", func() {
			g.It("returns nil when line is not an annotation", func() {
				annotation, err := ParseAnnotation("foo bar baz")
				g.Assert(annotation == nil).IsTrue("expected no annotation")
				g.Assert(err == nil).IsTrue("expected no error")
			})

			g.It("returns nil when annotation is missing centry prefix", func() {
				annotation, err := ParseAnnotation("foo/bar=baz")
				g.Assert(annotation == nil).IsTrue("expected no annotation")
				g.Assert(err == nil).IsTrue("expected no error")
			})

			g.It("returns nil when annotation is missing slash", func() {
				annotation, err := ParseAnnotation("centry_bar=baz")
				g.Assert(annotation == nil).IsTrue("expected no annotation")
				g.Assert(err == nil).IsTrue("expected no error")
			})

			g.It("returns nil when annotation is missing equals sign", func() {
				annotation, err := ParseAnnotation("centry/bar_baz")
				g.Assert(annotation == nil).IsTrue("expected no annotation")
				g.Assert(err == nil).IsTrue("expected no error")
			})

			g.It("returns error when annotation only has equals sign before slash", func() {
				annotation, err := ParseAnnotation("centry=bar/bar")
				g.Assert(annotation == nil).IsTrue("expected no annotation")
				g.Assert(err != nil).IsTrue("expected no error")
			})
		})

		g.Describe("ValidAnnotations", func() {
			g.It("returns annotation", func() {
				annotation, err := ParseAnnotation("centry/foo=bar")
				g.Assert(annotation != nil).IsTrue("expected annotation")
				g.Assert(err == nil).IsTrue("expected no error")
				g.Assert(annotation.Namespace).Equal("centry")
				g.Assert(annotation.Key).Equal("foo")
				g.Assert(annotation.Value).Equal("bar")
			})

			g.It("returns annotation when it contains multiple equal signs", func() {
				annotation, err := ParseAnnotation("centry/foo=bar=baz")
				g.Assert(annotation != nil).IsTrue("expected no annotation")
				g.Assert(err == nil).IsTrue("expected no error")
				g.Assert(annotation.Namespace).Equal("centry")
				g.Assert(annotation.Key).Equal("foo")
				g.Assert(annotation.Value).Equal("bar=baz")
			})

			g.It("returns annotation when it contains multiple slashes", func() {
				annotation, err := ParseAnnotation("centry/foo=bar/baz")
				g.Assert(annotation != nil).IsTrue("expected no annotation")
				g.Assert(err == nil).IsTrue("expected no error")
				g.Assert(annotation.Namespace).Equal("centry")
				g.Assert(annotation.Key).Equal("foo")
				g.Assert(annotation.Value).Equal("bar/baz")
			})

			g.It("returns annotation when namespace contains key[value]", func() {
				annotation, err := ParseAnnotation("centry.key[value]/foo=bar")
				g.Assert(annotation != nil).IsTrue("expected no annotation")
				g.Assert(err == nil).IsTrue("expected no error")
				g.Assert(annotation.Namespace).Equal("centry.key")
				g.Assert(annotation.NamespaceValues["key"]).Equal("value")
				g.Assert(annotation.Key).Equal("foo")
				g.Assert(annotation.Value).Equal("bar")
			})

			g.It("returns annotation when namespace contains multiple key[value]", func() {
				annotation, err := ParseAnnotation("centry.key1[value1].key2[value2]/foo=bar")
				g.Assert(annotation != nil).IsTrue("expected no annotation")
				g.Assert(err == nil).IsTrue("expected no error")
				g.Assert(annotation.Namespace).Equal("centry.key1.key2")
				g.Assert(annotation.NamespaceValues["key1"]).Equal("value1")
				g.Assert(annotation.NamespaceValues["key2"]).Equal("value2")
				g.Assert(annotation.Key).Equal("foo")
				g.Assert(annotation.Value).Equal("bar")
			})

			g.It("returns annotation when namespace contains key[value] where value has special character", func() {
				annotation, err := ParseAnnotation("centry.key1[value:1].key2[value_2]/foo=bar")
				g.Assert(annotation != nil).IsTrue("expected no annotation")
				g.Assert(err == nil).IsTrue("expected no error")
				g.Assert(annotation.Namespace).Equal("centry.key1.key2")
				g.Assert(annotation.NamespaceValues["key1"]).Equal("value:1")
				g.Assert(annotation.NamespaceValues["key2"]).Equal("value_2")
				g.Assert(annotation.Key).Equal("foo")
				g.Assert(annotation.Value).Equal("bar")
			})
		})
	})
}
