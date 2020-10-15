package test

import (
	"fmt"

	"github.com/franela/goblin"
)

// AssertError asserts that error is not nil
func AssertError(g *goblin.G, e error) {
	if e == nil {
		g.Fail("expected an error")
	}
}

// AssertNoError asserts that error is nil
func AssertNoError(g *goblin.G, e error) {
	if e != nil {
		g.Fail(fmt.Sprintf("expected no error but got, %v", e))
	}
}
