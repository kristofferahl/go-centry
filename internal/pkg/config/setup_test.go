package config

import (
	"os"
)

func init() {
	// Esuring the workdir is the root of the repo
	os.Chdir("../../../")
}
