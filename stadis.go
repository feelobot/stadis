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

	cli.NewApp().Run(os.Args)
	app := cli.NewApp()
	app.Name = "stadis"
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
	}
	app.Action = func(c *cli.Context) {
		for {
			info := getStats(c.String("redis-host"))
			gauges := parseGauges(info)
			sendGauges(gauges, c.String("statsd-host"))
			color.Yellow("-------------------")
			time.Sleep(5000 * time.Millisecond)
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

func sendStats() {
}
func sendGauges(gauges map[string]int64, statsdHost string) {
	color.Yellow("-------------------")
	color.White("GAUGES:")
	client, err := statsd.NewClient(statsdHost, "test-client")
	if err != nil {
		color.Red("ERROR: can't connect to statsd")
	}
	defer client.Close()
	for gauge, value := range gauges {
		client.Gauge(gauge, value, 1.0)
	}

}
func statify(info string) {
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
