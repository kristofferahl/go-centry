package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	api "github.com/kristofferahl/go-centry/internal/pkg/api"
	"github.com/kristofferahl/go-centry/internal/pkg/config"
	"github.com/kristofferahl/go-centry/internal/pkg/io"
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
		Log:       sc.Log,
		BasicAuth: configureBasicAuth(),
	})

	s.Router.HandleFunc("/", sc.indexHandler()).Methods("GET")
	s.Router.HandleFunc("/commands/", sc.executeHandler()).Methods("POST")

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

func configureBasicAuth() *api.BasicAuth {
	var auth *api.BasicAuth
	baUsername := os.Getenv("CENTRY_SERVE_USERNAME")
	baPassword := os.Getenv("CENTRY_SERVE_PASSWORD")

	if baUsername != "" && baPassword != "" {
		auth = &api.BasicAuth{
			Username: baUsername,
			Password: baPassword,
		}
	}

	return auth
}

func (sc *ServeCommand) indexHandler() func(w http.ResponseWriter, r *http.Request) {
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

func (sc *ServeCommand) executeHandler() func(w http.ResponseWriter, r *http.Request) {
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
		args = append(args, sc.Manifest.Path)
		args = append(args, strings.Fields(body.Args)...)

		// Build
		io, buf := io.BufferedCombined()
		context := NewContext(API, io)

		context.commandEnabledFunc = func(cmd config.Command) bool {
			if cmd.Annotations == nil || cmd.Annotations[config.CommandAnnotationAPIServe] != config.TrueString {
				return false
			}

			return true
		}

		context.optionEnabledFunc = func(opt config.Option) bool {
			if opt.Annotations == nil || opt.Annotations[config.CommandAnnotationAPIServe] != config.TrueString {
				return false
			}

			return true
		}

		runtime, err := NewRuntime(args, context)
		if err != nil {
			response.Centry = fmt.Sprintf("%s %s", context.manifest.Config.Name, context.manifest.Config.Version)
			response.Result = fmt.Sprintf("Unable to create runtime %v", err)
			response.ExitCode = 1
			sc.Log.Error(response.Result)
		} else {
			// Run
			exitCode := runtime.Execute()

			response.Centry = fmt.Sprintf("%s %s", context.manifest.Config.Name, context.manifest.Config.Version)
			response.Result = buf.String()
			response.ExitCode = exitCode
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		js, err := json.Marshal(response)
		if err == nil {
			w.Write(js)
		}
	}
}
