package main

import (
	. "aiven-prometheus-exporter/internal/pkg"
	"flag"
	"github.com/aiven/aiven-go-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {

	flag.BoolVar(&debugEnabled, "debug", false, "Enable debug logging")
	flag.Parse()

	if debugEnabled {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debugging enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	aiven_token, set := os.LookupEnv("AIVEN_API_TOKEN")
	if set == false {
		log.Fatal("No Aiven token found in environment. Please export the token as environment variable")
	}

	aivenClient, err := aiven.NewTokenClient(aiven_token, "")
	if err != nil {
		log.Fatal(err)
	}

	exporter := AivenCollector{Client: aivenClient}
	r := prometheus.NewRegistry()
	r.MustRegister(exporter)

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":2112", nil))

}

var (
	debugEnabled = false
)
