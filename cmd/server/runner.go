package main

import (
	"fmt"
	"github.com/quan-to/slog"
	"github.com/racerxdl/radioserver"
	"github.com/racerxdl/radioserver/server"
	"github.com/racerxdl/segdsp/dsp"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"syscall"
)

var log = slog.Scope("RadioServer")

func main() {
	addr, freq, gain := getConfig()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Got panic", r)
			debug.PrintStack()
			os.Exit(255)
		}
	}()

	log.Info("Initializing Frontend")
	frontend := createFrontend()
	frontend.Init()
	defer frontend.Destroy()

	frontend.SetCenterFrequency(uint32(freq))
	if gain > 0 {
		frontend.SetGain(uint8(gain))
	}
	frontend.Start()

	defer frontend.Stop()

	log.Info("Protocol Version: %s", radioserver.ServerVersion.AsString())
	log.Info("SIMD Mode: %s", dsp.GetSIMDMode())

	srv := server.MakeRadioServer(frontend)
	err := srv.Listen(addr)
	if err != nil {
		log.Error("Error listening: %s", err)
	}
	stop := make(chan bool, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c
		log.Info("Got SIGTERM! Closing it")
		stop <- true
	}()

	<-stop

	srv.Stop()
	log.Info("Done")
}
