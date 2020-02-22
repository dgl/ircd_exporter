package irc

import (
	"time"
)

type StatsRequest struct {
	Local    bool
	Timeout  time.Duration
	response chan StatsResponse
}

type StatsResponse struct {
	Timeout  bool
	Servers  map[string]*ServerStats
	Channels int
}

type ServerStats struct {
	Up, done                  bool
	RequestTime, ResponseTime time.Time
	Distance, Users           int
}
