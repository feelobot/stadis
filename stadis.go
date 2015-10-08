package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/quipo/statsd"
	"gopkg.in/redis.v3"
	"os"
	"regexp"
	"strconv"
	"time"
)

var prefix = os.Getenv("HOSTNAME")
var statsdclient = statsd.NewStatsdClient("localhost:8125", prefix)

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
		info := getStats(c.String("redis-host"))
		gauges := parseGauges(info)
		sendGauges(gauges)
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
		fmt.Println("Error connecting to redis") //possibly send to statsd also
	}
	//fmt.Println(info)
	return info
}

func sendStats() {
}
func sendGauges(gauges map[string]int64) {
	statsdclient.CreateSocket()
	interval := time.Second * 5 // aggregate stats and flush every 2 seconds
	stats := statsd.NewStatsdBuffer(interval, statsdclient)
	defer stats.Close()

	for k, v := range gauges {
		fmt.Println("k:", k, "v:", v)
		stats.Gauge("k", v)
	}

}
func statify(info string) {
}

func parseGauges(info string) map[string]int64 {
	var gauges_with_values map[string]int64
	gauges := []string{
		"blocked_clients",
		"connected_clients",
		"instantaneous_ops_per_sec",
		"latest_fork_usec",
		"mem_fragmentation_ratio",
		"migrate_cached_sockets",
		"pubsub_channels",
		"pubsub_patterns",
		"uptime_in_seconds",
		"used_memory",
		"used_memory_lua",
		"used_memory_peak",
		"used_memory_rss",
	}
	for _, gauge := range gauges {
		r, _ := regexp.Compile(fmt.Sprint(gauge, ":([0-9]*)"))
		value := r.FindStringSubmatch(info)[1]
		fmt.Println(fmt.Sprint(gauge, ": ", value))
		v, _ := strconv.ParseInt(value, 10, 64)
		gauges_with_values[gauge] = v
	}
	return gauges_with_values
}
