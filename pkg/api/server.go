package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Config defines server configuration
type Config struct {
	Log *logrus.Entry
}

// Server holds the configuration, router and http server
type Server struct {
	Config Config
	Router *mux.Router
	Server *http.Server
}

// NewServer creates a new HTTP server
func NewServer(config Config) Server {
	s := Server{
		Server: nil,
		Config: config,
		Router: mux.NewRouter(),
	}

	s.Router.Use(loggingMiddleware(config))

	s.Server = &http.Server{
		Handler:      s.Router,
		Addr:         ":8113",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	config.Log.WithFields(logrus.Fields{
		"service": "HTTP-Server",
		"address": s.Server.Addr,
	}).Infof("listening on %s", s.Server.Addr)

	return s
}

// RunAndBlock starts the HTTP server and blocks
func (s Server) RunAndBlock() error {
	return s.Server.ListenAndServe()
}

func loggingMiddleware(config Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			config.Log.WithFields(logrus.Fields{
				"service":   "HTTP-Server",
				"method":    r.Method,
				"path":      r.RequestURI,
				"direction": "incoming",
			}).Debugf("HTTP %s %s", r.Method, r.RequestURI)

			scrw := &statusCodeResponseWriter{w, http.StatusOK}
			next.ServeHTTP(scrw, r)

			config.Log.WithFields(logrus.Fields{
				"service":   "HTTP-Server",
				"status":    scrw.statusCode,
				"direction": "outgoing",
			}).Debugf("HTTP %s %s - %d %s", r.Method, r.RequestURI, scrw.statusCode, http.StatusText(scrw.statusCode))
		})
	}
}

type statusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *statusCodeResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
