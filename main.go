package main

import (
	"os"

	"github.com/kristofferahl/go-centry/pkg/centry"
)

func main() {
	os.Exit(centry.RunOnce(os.Args))
}
