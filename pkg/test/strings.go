package test

import (
	"fmt"
	"strings"

	"github.com/franela/goblin"
)

// AssertKeyValueExists asserts the given key and value is present on one of the lines given as input
func AssertKeyValueExists(g *goblin.G, key, value, input string) {
	found := false
	lines := strings.Split(input, "\n")
	for _, l := range lines {
		parts := strings.Split(l, "=")
		k := parts[0]

		var v string
		if len(parts) > 1 {
			v = parts[1]
		}

		if k == key {
			found = true
			g.Assert(v == value).IsTrue(fmt.Sprintf("wrong expected value for key \"%s\" expected=%s actual=%s", key, value, v))
		}
	}

	if !found {
		g.Fail(fmt.Sprintf("\"%s\" key not found in input", key))
	}
}
