package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"text/template"

	"net/http"
	"os"
	"time"

	"github.com/MihaiBlebea/go-wallet/bot"
	"github.com/MihaiBlebea/go-wallet/domain"
	"github.com/MihaiBlebea/go-wallet/server/handler"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//go:embed templates/*
var templates embed.FS

//go:embed static/*
var static embed.FS

const prefix = "/api/v1/"

func NewServer(wallet domain.Wallet, bot bot.Service, logger *logrus.Logger) {
	r := mux.NewRouter()

	api := r.PathPrefix(prefix).
		Subrouter()

	api.Handle("/health-check", handler.HealthEndpoint()).
		Methods("GET")

	api.Handle("/webhook", handler.WebhookEndpoint(wallet, bot)).
		Methods("POST")

	api.Handle("/report", handler.ReportEndpoint(wallet)).
		Methods("GET")

	tmpl := template.Must(template.ParseFS(templates, "templates/*"))

	r.Handle("/report", handler.WebappEndpoint(tmpl, wallet)).
		Methods("GET")

	statics, err := fs.Sub(static, "static")
	if err != nil {
		log.Panic(err)
	}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(statics))))

	r.Use(loggerMiddleware(logger))

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func loggerMiddleware(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("Incoming %s request %s", r.Method, r.URL.Path))
			next.ServeHTTP(w, r)
		})
	}
}
