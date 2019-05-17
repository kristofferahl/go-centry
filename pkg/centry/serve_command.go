package centry

import (
	"encoding/json"
	"net/http"
	"strings"

	api "github.com/kristofferahl/go-centry/pkg/api"
	"github.com/kristofferahl/go-centry/pkg/config"
	"github.com/kristofferahl/go-centry/pkg/io"
	"github.com/sirupsen/logrus"
)

// ServeCommand is a Command implementation that applies stuff
type ServeCommand struct {
	Manifest *config.Manifest
	Log      *logrus.Entry
}

// Run starts an HTTP server and blocks
func (sc *ServeCommand) Run(args []string) int {
	sc.Log.Debugf("Serving HTTP api")

	s := api.NewServer(api.Config{
		Log: sc.Log,
	})

	s.Router.HandleFunc("/", indexHandler(sc.Manifest)).Methods("GET")
	s.Router.HandleFunc("/commands/", executeHandler(sc.Manifest)).Methods("POST")

	err := s.RunAndBlock()
	if err != nil {
		return 1
	}

	return 0
}

// Help returns the help text of the ServeCommand
func (sc *ServeCommand) Help() string {
	return "No help here..."
}

// Synopsis returns the synopsis of the ServeCommand
func (sc *ServeCommand) Synopsis() string {
	return "Exposes commands over HTTP"
}

func indexHandler(manifest *config.Manifest) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusOK
		response := api.IndexResponse{}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		js, err := json.Marshal(response)
		if err == nil {
			w.Write(js)
		}
	}
}

func executeHandler(manifest *config.Manifest) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := http.StatusOK
		response := api.ExecuteResponse{}

		var body api.ExecuteRequest

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		if err != nil {
			statusCode = http.StatusBadRequest
		}

		args := []string{}
		args = append(args, manifest.Path)
		args = append(args, strings.Fields(body.Args)...)

		// Build
		io, buf := io.BufferedCombined()
		context := NewContext(API, io)

		context.commandEnabled = func(cmd config.Command) bool {
			if cmd.Annotations == nil || cmd.Annotations[config.APIServeAnnotation] != "true" {
				return false
			}

			return true
		}

		runtime := Create(args, context)

		// Run
		exitCode := runtime.Execute()

		response.Result = buf.String()
		response.ExitCode = exitCode

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		js, err := json.Marshal(response)
		if err == nil {
			w.Write(js)
		}
	}
}
