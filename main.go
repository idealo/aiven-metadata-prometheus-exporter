package main

import (
	. "aiven-prometheus-exporter/internal/pkg"
	"flag"
	"github.com/aiven/aiven-go-client"
	"github.com/go-co-op/gocron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func main() {

	flag.BoolVar(&debugEnabled, "debug", false, "Enable debug logging")
	flag.Parse()

	setupLogging()

	aiven_token, set := os.LookupEnv("AIVEN_API_TOKEN")
	if set == false {
		log.Fatal("No Aiven token found in environment. Please export the token as 'AIVEN_API_TOKEN' environment variable")
	}

	aivenClient, err := aiven.NewTokenClient(aiven_token, "")
	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("Starting the Aiven Prometheus Exporter")

	exporter := AivenCollector{Client: aivenClient}
	r := prometheus.NewRegistry()
	r.MustRegister(exporter)

	scheduler := gocron.NewScheduler(time.UTC)
	interval := "5m"
	_, err = scheduler.Every(interval).Do(func() { exporter.CollectAsync() })
	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("Scheduler set up. Metrics will be refreshed every", interval)
	scheduler.StartAsync()

	log.Infoln("Listening on port 2112 and /metrics")
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":2112", nil))

}

func setupLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	if debugEnabled {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debugging enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

var (
	debugEnabled = false
)
