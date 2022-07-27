package main

import (
	. "aiven-metadata-prometheus-exporter/internal/pkg"
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
	flag.StringVar(&interval, "scrape-interval", "5m", "Aiven API scrape interval. Defaults to 5m")
	flag.Parse()

	setupLogging()

	aivenToken, set := os.LookupEnv("AIVEN_API_TOKEN")
	if set == false {
		log.Fatal("No Aiven token found in environment. Please export the token as 'AIVEN_API_TOKEN' environment variable")
	}

	aivenClient, err := aiven.NewTokenClient(aivenToken, "")
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Starting Aiven Metadata Prometheus Exporter")

	collector := AivenCollector{AivenClient: aivenClient}
	r := prometheus.NewRegistry()
	r.MustRegister(collector)

	scheduler := gocron.NewScheduler(time.UTC)
	_, err = scheduler.Every(interval).Do(func() { collector.CollectAsync() })
	if err != nil {
		log.Fatalln(err)
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
	interval     string
)
