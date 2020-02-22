package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bitnami-labs/flagenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/dgl/prometheus-ircd-user-exporter/irc"
)

var (
	ircOptions irc.Options
	flagListen = flag.String("listen", ":9678", "[host]:port to serve HTTP on, for metrics collection")
)

func init() {
	irc.Flags("irc.", &ircOptions)
	flagenv.SetFlagsFromEnv("PIE", flag.CommandLine)
}

func main() {
	flag.Parse()

	// Connect to IRC
	client := irc.NewClient(ircOptions)
	go client.Start()

	prometheus.MustRegister(&Exporter{client})

	http.HandleFunc("/-/healthy", ok)
	http.HandleFunc("/-/ready", ok)

	http.Handle("/metrics", promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer,
		promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{})))

	log.Fatal(http.ListenAndServe(*flagListen, nil))
}

func ok(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
