package main

import (
	"flag"
	"time"
)

var flags struct {
	listenPort      int
	monitorInterval time.Duration
}

func init() {
	flag.IntVar(&flags.listenPort, "p", 8080, "Listen port.")
	flag.DurationVar(&flags.monitorInterval, "monitor-interval", 10*time.Second, "Log for traffic cache status will print out per monitor-interval.")
}
