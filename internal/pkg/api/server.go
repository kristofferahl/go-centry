package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Config defines server configuration
type Config struct {
	Log       *logrus.Entry
	BasicAuth *BasicAuth
}

// BasicAuth defines basic auth configuration
type BasicAuth struct {
	Username string
	Password string
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
	s.Router.Use(basicAuthMiddleware(config))

	s.Server = &http.Server{
		Handler:      s.Router,
		Addr:         ":8113",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	auth := "none"
	if config.BasicAuth != nil {
		auth = "basic"
	}

	l := config.Log.WithFields(logrus.Fields{
		"service": "HTTP-Server",
		"address": s.Server.Addr,
		"auth":    auth,
	})

	if auth != "basic" {
		l.Warnf("listening on %s", s.Server.Addr)
	} else {
		l.Infof("listening on %s", s.Server.Addr)
	}

	return s
}

// RunAndBlock starts the HTTP server and blocks
func (s Server) RunAndBlock() error {
	return s.Server.ListenAndServe()
}

func basicAuthMiddleware(config Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ba := config.BasicAuth
			if ba == nil {
				next.ServeHTTP(w, r)
				return
			}

			user, pass, _ := r.BasicAuth()

			if (*ba).Username != user || (*ba).Password != pass {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				config.Log.WithFields(logrus.Fields{
					"service":    "HTTP-Server",
					"middleware": "basic-auth",
				}).Debugf("authentication failed")
				return
			}

			config.Log.WithFields(logrus.Fields{
				"service":    "HTTP-Server",
				"middleware": "basic-auth",
			}).Debugf("successfully authenticated")

			next.ServeHTTP(w, r)
		})
	}
}

func loggingMiddleware(config Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			config.Log.WithFields(logrus.Fields{
				"service":    "HTTP-Server",
				"middleware": "logging",
				"method":     r.Method,
				"path":       r.RequestURI,
				"direction":  "incoming",
			}).Debugf("HTTP %s %s", r.Method, r.RequestURI)

			scrw := &statusCodeResponseWriter{w, http.StatusOK}
			next.ServeHTTP(scrw, r)

			config.Log.WithFields(logrus.Fields{
				"service":    "HTTP-Server",
				"middleware": "logging",
				"status":     scrw.statusCode,
				"direction":  "outgoing",
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
