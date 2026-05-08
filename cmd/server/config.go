package main

import (
	"flag"
	"os"
	"strconv"
)

var listenAddr = flag.String("listen", ":4050", "Listen address (host:port)")
var frequency = flag.Uint("frequency", 106300000, "Initial center frequency in Hz")
var gain = flag.Uint("gain", 0, "Initial gain")
var cpuprofile = flag.String("cpuprofile", "", "Write CPU profile to file")

func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envUint(key string, fallback uint) uint {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return fallback
		}
		return uint(n)
	}
	return fallback
}

func getConfig() (addr string, freq uint, g uint) {
	flag.Parse()
	addr = envString("LISTEN_ADDR", *listenAddr)
	freq = envUint("CENTER_FREQUENCY", *frequency)
	g = envUint("GAIN", *gain)
	return
}
