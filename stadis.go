package main

import (
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"gopkg.in/redis.v3"
	"os"
	"regexp"
	"strconv"
	"time"
)

func main() {
	// init
	app := cli.NewApp()
	app.Name = "stadis"
	app.Version = "0.0.1"
	app.Usage = "get redis info and submit to statsd"
	app.HideHelp = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "redis-host, r",
			Value: "localhost:6379",
			Usage: "host:port of redis servier",
		},
		cli.StringFlag{
			Name:  "statsd-host, s",
			Value: "localhost:8125",
			Usage: "host:port of statsd servier",
		},
		cli.StringFlag{
			Name:   "prefix,p",
			Usage:  "host:port of redis servier",
			EnvVar: "HOSTNAME",
		},
		cli.StringFlag{
			Name:  "interval,i",
			Usage: "time in milliseconds to periodically check redis",
			Value: "5000",
		},
	}
	app.Action = func(c *cli.Context) {
		for {
			info := getStats(c.String("redis-host"))
			gauges := parseGauges(info)
			counters := parseCounters(info)
			sendStats(c.String("statsd-host"), c.String("prefix"), gauges, counters)
			color.Yellow("-------------------")
			interval, _ := strconv.ParseInt(c.String("interval"), 10, 64)
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}
	app.Run(os.Args)
}

func getStats(addrs string) string {
	client := redis.NewClient(&redis.Options{
		Addr:     addrs,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	info, err := client.Info().Result()
	if err != nil {
		color.Red("ERROR: can't connect to redis") //possibly send to statsd also
	}
	// DEBUG: fmt.Println(info)
	return info
}

func sendStats(statsdHost string, prefix string, gauges map[string]int64, counters map[string]int64) {
	client, err := statsd.NewClient(statsdHost, "")
	if err != nil {
		color.Red("ERROR: can't connect to statsd")
	}
	defer client.Close()
	sendGauges(client, gauges)
	sendCounters(client, counters)
}

func sendGauges(client statsd.Statter, gauges map[string]int64) {
	for gauge, value := range gauges {
		client.Gauge(gauge, value, 1.0)
	}
}

func sendCounters(client statsd.Statter, counters map[string]int64) {
	for counter, value := range counters {
		client.Inc(counter, value, 1.0)
	}
}

func parseGauges(info string) map[string]int64 {
	gauges_with_values := map[string]int64{
		"blocked_clients":           0,
		"connected_clients":         0,
		"instantaneous_ops_per_sec": 0,
		"latest_fork_usec":          0,
		"mem_fragmentation_ratio":   0,
		"migrate_cached_sockets":    0,
		"pubsub_channels":           0,
		"pubsub_patterns":           0,
		"uptime_in_seconds":         0,
		"used_memory":               0,
		"used_memory_lua":           0,
		"used_memory_peak":          0,
		"used_memory_rss":           0,
	}
	color.Yellow("-------------------")
	color.White("GAUGES:")
	for gauge, _ := range gauges_with_values {
		r, _ := regexp.Compile(fmt.Sprint(gauge, ":([0-9]*)"))
		matches := r.FindStringSubmatch(info)
		if matches == nil {
			color.Red(fmt.Sprint("ERROR: ", gauge, "is not displayed in redis info"))
		} else {
			value := matches[len(matches)-1]
			color.Cyan(fmt.Sprint(gauge, ": ", value))
			v, _ := strconv.ParseInt(value, 10, 64)
			gauges_with_values[gauge] = v
		}
	}
	return gauges_with_values
}

func parseCounters(info string) map[string]int64 {
	counters := map[string]int64{
		"evicted_keys":               0,
		"expired_keys":               0,
		"keyspace_hits":              0,
		"keyspace_misses":            0,
		"rejected_connections":       0,
		"sync_full":                  0,
		"sync_partial_err":           0,
		"sync_partial_ok":            0,
		"total_commands_processed":   0,
		"total_connections_received": 0,
	}
	color.Yellow("-------------------")
	color.White("COUNTERS:")
	for counter, _ := range counters {
		r, _ := regexp.Compile(fmt.Sprint(counter, ":([0-9]*)"))
		matches := r.FindStringSubmatch(info)
		if matches == nil {
			color.Red(fmt.Sprint("ERROR: ", counter, "is not displayed in redis info"))
		} else {
			value := matches[len(matches)-1]
			color.Cyan(fmt.Sprint(counter, ": ", value))
			v, _ := strconv.ParseInt(value, 10, 64)
			counters[counter] = v
		}
	}
	return counters
}
