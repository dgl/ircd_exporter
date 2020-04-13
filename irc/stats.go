package irc

import (
	"time"
)

type StatsRequest struct {
	Local         bool
	Timeout       time.Duration
	IgnoreServers []string
	Nicks         []string
	response      chan StatsResponse
}

type StatsResponse struct {
	Timeout  bool
	Servers  map[string]*ServerStats
	Channels int
	Nicks    map[string]bool
}

type ServerStats struct {
	Up, done                  bool
	RequestTime, ResponseTime time.Time
	Distance, Users           int
}
