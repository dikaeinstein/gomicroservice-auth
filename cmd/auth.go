package main

import (
	"net/http"
	"os"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/dikaeinstein/gomicroservice-auth/config"
	"github.com/dikaeinstein/gomicroservice-auth/handler"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	const address = ":8081"
	logger := &log.Logger{
		Out:       os.Stdout,
		Formatter: new(log.TextFormatter),
		Level:     log.DebugLevel,
	}

	cfg := config.New()
	c, err := statsd.New(cfg.DogStatsD)
	if err != nil {
		log.Fatal(err)
	}

	c.Namespace = "gomicroservice.auth."

	jwt := handler.NewJWT(logger, c, cfg)
	health := handler.NewHealth(logger, c)

	r := mux.NewRouter()
	r.HandleFunc("/login", jwt.Login).Methods(http.MethodPost)
	r.HandleFunc("/healthz", health.Get).Methods(http.MethodGet)

	logger.WithField("service", "auth").
		Infof("Starting server, listening on %s", address)
	log.WithField("service", "auth").
		Fatal(http.ListenAndServe(address, r))
}
