# ircd_exporter

Export user counts and related metadata from an IRC network to Prometheus.

This is an IRC client, that stays connected, then collects some metadata from
the server when queried by `/metrics`.

This should work with most IRC servers, but has so far only been tested on Charybdis.

## Using

### Manual build

```shell
go install github.com/dgl/ircd_exporter/cmd/ircd_exporter@latest
```

This will give you a *ircd_exporter* to run, in your Go bin directory:

```shell
$(go env GOPATH)/bin/ircd_exporter
```

### Docker

There is a [docker image](https://github.com/dgl/ircd_exporter/pkgs/container/ircd_exporter).
You can use this with the example [docker-compose.yml](docker-compose.yml) or run it directly with Docker:

```cli
docker run --restart=unless-stopped -p 9678:9678 -e PIE_IRC.SERVER=irc.example.com:6667 ghcr.io/dgl/ircd_exporter
```

## Options

Server to connect to:

```
  -irc.nick string
    	Nickname to use (default "promexp")
  -irc.oper string
    	Username to use for /OPER (optional)
  -irc.oper-password string
    	Password to use for /OPER (optional)
  -irc.password string
    	Password to use when connecting to the server (optional)
  -irc.server string
    	Server to connect to, host:port (default "localhost:6667")
```

Where to listen for /metrics requests:

```
  -listen string
    	[host]:port to serve HTTP on, for metrics collection (default ":9678")
```

Stats to collect:

```
  -stats.ignore string
        Servers to ignore for stats (comma separated, e.g. some services servers
        don't support the LUSERS command).
  -stats.local-only
        Only get stats from the local server. Default is to run /LINKS and
        record the LUSERS output from them all.
  -stats.nicks
        List of nicknames to check for ISON status (comma separated).
  -stats.timeout duration
        How long to wait before for stats reply before considering a server
        down. (default 9s)
```

Environment variables are also supported for all flags (prefix with `PIE_` and
replace `-` with `_`), but useful for passwords in particular, e.g.:
`PIE_IRC.OPER_PASSWORD`. (Many servers need operator privileges to see full
information in /LINKS, some restrict remote commands, if you want the timings to
be accurate you need flood controls to not be applied, etc.
So making this an operator is recommended.)

If you don't want to make this an operator using `--stats.local-only` and
running one per IRC server is recommended. Note `--stats.local-only` only
reports the server it is directly connected to in metrics such as `irc_up`, etc.
However `irc_distance` includes all servers (based on /LINKS output), see
[example.yaml](example.yaml) for a way to use this metric usefully.

## Metrics exported

```
# HELP irc_connected Is the exporter connected to the server?
# TYPE irc_connected gauge
irc_connected 1

# HELP irc_channels Number of channels created in the IRC network.
# TYPE irc_channels gauge
irc_channels 42

# HELP irc_distance Number of hops this server is in the IRC network from the server where the exporter is running.
# TYPE irc_distance gauge
irc_distance{server="local.example.org"} 0
irc_distance{server="remote.example.org"} 1

# HELP irc_latency_seconds Latency of the request to this server from where the exporter is running.
# TYPE irc_latency_seconds gauge
irc_latency_seconds{server="local.example.org"} 0.015440223

# HELP irc_up Was the last query of each server successful.
# TYPE irc_up gauge
irc_up{server="local.example.org"} 1
irc_up{server="remote.example.org"} 0

# HELP irc_users Number of users on this IRC server.
# TYPE irc_users gauge
irc_users{server="local.example.org"} 746

# HELP irc_ison Whether specified nicknames are online or not.
# TYPE irc_ison gauge
irc_ison{nick="fred"} 0
irc_ison{nick="nickserv"} 1
```

## Prometheus config

Just scrape this, e.g.:

```yaml
scrape_configs:
  - job_name: 'ircd-exporter'
    static_configs:
      - targets: ['localhost:9678']
```

See [example.yaml](example.yaml) for some ways to use the metrics.

## Alternatives

* https://github.com/wobscale/prometheus-irc-exporter
* https://github.com/wikimedia/operations-debs-prometheus-ircd-exporter/blob/master/prometheus-ircd-exporter

These both connect when scraped, I didn't want to do that (you can use the
blackbox exporter to test a TCP or TLS connection... [see
example](https://github.com/prometheus/blackbox_exporter/blob/bf3e7fbbec35ce1b5ffd2c7abdf3ebc9ec4bc975/blackbox.yml#L23)).
