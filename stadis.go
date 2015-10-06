package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/quipo/statsd"
	"gopkg.in/redis.v3"
	"os"
	"time"
)

func main() {
	// init
	cli.NewApp().Run(os.Args)
	app := cli.NewApp()
	app.Name = "stadis"
	app.Usage = "get redis info and submit to statsd"
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
		getStats(c.String("redis-host"))
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
	fmt.Println(info)
	return info
}

func sendStats() {
	prefix := os.Getenv("FOO")
	statsdclient := statsd.NewStatsdClient("localhost:8125", prefix)
	statsdclient.CreateSocket()
	interval := time.Second * 5 // aggregate stats and flush every 2 seconds
	stats := statsd.NewStatsdBuffer(interval, statsdclient)
	defer stats.Close()

	// not buffered: send immediately
	statsdclient.Incr("mymetric", 4)

	// buffered: aggregate in memory before flushing
	stats.Incr("mymetric", 1)
	stats.Incr("mymetric", 3)
}
