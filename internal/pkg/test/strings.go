package test

import (
	"fmt"
	"strings"

	"github.com/franela/goblin"
)

const tmplExpectedKey string = `expected key "%s" to be found in input

INPUT
----------------------------------------------------
%s
----------------------------------------------------`

const tmplExpectedContains string = `

EXPECTED THIS
----------------------------------------------------
%s
----------------------------------------------------

TO BE FOUND IN
----------------------------------------------------
%s
----------------------------------------------------`

const tmplExpectedValue string = `expected value "%s" for key "%s" but found "%s"`

// TODO: Remove
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
		g.Fail(fmt.Sprintf("\"%s\" key not found in input:\n\n%s", key, input))
	}
}

// AssertStringHasKeyValue asserts the expected string is found in within the input
func AssertStringHasKeyValue(g *goblin.G, s, key, value string) {
	found := false
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		parts := strings.Split(l, "=")
		k := parts[0]

		var v string
		if len(parts) > 1 {
			v = parts[1]
		}

		if k == key {
			found = true
			if v != value {
				g.Fail(fmt.Sprintf(tmplExpectedValue, value, key, v))
			}
		}
	}

	if !found {
		g.Fail(fmt.Sprintf(tmplExpectedKey, key, s))
	}
}

// AssertStringContains asserts the expected string is found in within the input
func AssertStringContains(g *goblin.G, s, substring string) {
	s = strings.TrimSpace(s)
	substring = strings.TrimSpace(substring)
	msg := fmt.Sprintf(tmplExpectedContains, substring, s)
	g.Assert(strings.Contains(s, substring)).IsTrue(msg)
}
