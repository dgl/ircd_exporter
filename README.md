# prometheus-ircd-user-exporter

Export user counts and related metadata from an IRC network to Prometheus.

This is an IRC client, that stays connected, then collects some metadata from
the server when queried by `/metrics`.

## Options

* --irc.server
* --irc.password
* --irc.nick
* --irc.oper
* --irc.oper-password

* --stats.local-only
* --stats.timeout

Environment variables are also supported all flags, but useful for passwords in
particular, e.g.: `IRC.OPER_PASSWORD`. (Many servers need operator privileges to
see full information in /LINKS, some restrict remote commands, if you want the
timings to be accurate you need flood controls to not be applied, etc. So making
this an operator is recommended.)

If you don't want to make it an operator using `--stats.local-only` and running
one per IRC server is recommended.

## Metrics exported

```
# HELP irc_channels Number of channels created in the IRC network.
# TYPE irc_channels gauge
irc_channels 42

# HELP irc_distance Number of hops this server is in the IRC network from the server where the exporter is running.
# TYPE irc_distance gauge
irc_distance{server="local.example.org"} 0
irc_distance{server="remote.example.org"} 1

# HELP irc_latency Latency of the request to this server from where the exporter is running.
# TYPE irc_latency gauge
irc_latency{server="local.example.org"} 0.015440223

# HELP irc_up Was the last query of each server successful.
# TYPE irc_up gauge
irc_up{server="local.example.org"} 1
irc_up{server="remote.example.org"} 0

# HELP irc_users Number of users on this IRC server.
# TYPE irc_users gauge
irc_users{server="local.example.org"} 746
```

## Alternatives

* https://github.com/wobscale/prometheus-irc-exporter/blob/master/main.go
* https://github.com/wikimedia/operations-debs-prometheus-ircd-exporter/blob/master/prometheus-ircd-exporter

These both connect when scraped, I didn't want to do that (you can use the
blackbox exporter to test a TCP or TLS connection...).
