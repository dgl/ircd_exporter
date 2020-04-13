package main

import (
	"flag"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/dgl/ircd_exporter/irc"
)

var (
	flagStatsLocal   = flag.Bool("stats.local-only", false, "Only get stats from the local server.")
	flagStatsTimeout = flag.Duration("stats.timeout", 9*time.Second, "How long to wait before for stats reply before considering a server down.")
	flagStatsIgnore  = flag.String("stats.ignore", "", "Servers to ignore for stats (comma separated, e.g. some services servers don't support the LUSERS command).")
	flagStatsNicks   = flag.String("stats.nicks", "", "List of nicknames to check for ISON status (comma separated).")
)

const (
	namespace = "irc"
)

var (
	connected = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "connected"),
		"Is the exporter connected to the server?",
		nil, nil,
	)
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of each server successful.",
		[]string{"server"}, nil,
	)
	distance = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "distance"),
		"Number of hops this server is in the IRC network from the server where the exporter is running.",
		[]string{"server"}, nil,
	)
	latency = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "latency_seconds"),
		"Latency of the request to this server from where the exporter is running.",
		[]string{"server"}, nil,
	)
	users = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "users"),
		"Number of users on this IRC server.",
		[]string{"server"}, nil,
	)
	channels = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "channels"),
		"Number of channels created in the IRC network.",
		nil, nil,
	)
	ison = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ison"),
		"Whether specified nicknames are online or not.",
		[]string{"nick"}, nil,
	)

	boolToFloat = map[bool]float64{
		false: 0.0,
		true:  1.0,
	}
)

type Exporter struct {
	client *irc.Client
}

// Describe describes all the metrics ever exported by the IRC exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- connected
	ch <- up
	ch <- distance
	ch <- latency
	ch <- users
	ch <- channels
	ch <- ison
}

// Collect gets stats from IRC and returns them as Prometheus metrics. It
// implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ignore := strings.Split(*flagStatsIgnore, ",")
	if len(ignore) == 1 && ignore[0] == "" {
		ignore = []string{}
	}
	nicks := strings.Split(*flagStatsNicks, ",")
	if len(nicks) == 1 && nicks[0] == "" {
		nicks = []string{}
	}
	res := e.client.Stats(irc.StatsRequest{
		Local:         *flagStatsLocal,
		Timeout:       *flagStatsTimeout,
		IgnoreServers: ignore,
		Nicks:         nicks,
	})

	ch <- prometheus.MustNewConstMetric(
		connected, prometheus.GaugeValue, boolToFloat[e.client.Server != ""])

	_, ok := res.Servers[e.client.Server]
	if res.Timeout && !ok {
		// Timeout, no data at all
		if e.client.Server != "" {
			ch <- prometheus.MustNewConstMetric(
				up, prometheus.GaugeValue, 0.0, e.client.Server)
		}
	} else {
		// Global state
		ch <- prometheus.MustNewConstMetric(
			channels, prometheus.GaugeValue, float64(res.Channels))

		for nick, nickIson := range res.Nicks {
			ch <- prometheus.MustNewConstMetric(
				ison, prometheus.GaugeValue, boolToFloat[nickIson], nick)
		}

		// Per server state
		for server, stats := range res.Servers {
			ch <- prometheus.MustNewConstMetric(
				distance, prometheus.GaugeValue, float64(stats.Distance), server)

			if *flagStatsLocal && e.client.Server != server {
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				up, prometheus.GaugeValue, boolToFloat[stats.Up], server)

			if stats.Up {
				ch <- prometheus.MustNewConstMetric(
					users, prometheus.GaugeValue, float64(stats.Users), server)

				ch <- prometheus.MustNewConstMetric(
					latency, prometheus.GaugeValue, float64(stats.ResponseTime.Sub(stats.RequestTime))/float64(time.Second), server)
			}
		}
	}
}
