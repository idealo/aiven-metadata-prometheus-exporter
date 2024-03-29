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

	collector := AivenCollector.Init(AivenCollector{}, aivenClient)
	r := prometheus.NewRegistry()
	r.MustRegister(collector)

	scheduler := gocron.NewScheduler(time.UTC)
	_, err = scheduler.Every(*interval).Do(func() { collector.CollectScheduled() })
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Scheduler set up. Metrics will be refreshed every", *interval)
	scheduler.StartAsync()

	log.Infoln("Listening on port", *listenAddress, "and", *metricsPath)
	http.Handle(*metricsPath, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}

func setupLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	if *debugEnabled {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debugging enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

var (
	debugEnabled  = flag.Bool("debug", false, "Enable debug logging")
	interval      = flag.String("collect-interval", "5m", "In which interval shall the metrics collected from the Aiven API")
	listenAddress = flag.String("listen-address", ":2112", "Address to listen on for telemetry")
	metricsPath   = flag.String("telemetry-path", "/metrics", "Path under which to expose metrics")
)
