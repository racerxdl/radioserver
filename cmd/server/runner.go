package main

import (
	"flag"
	"fmt"
	"github.com/quan-to/slog"
	"github.com/racerxdl/radioserver"
	"github.com/racerxdl/radioserver/frontends"
	"github.com/racerxdl/radioserver/server"
	"github.com/racerxdl/segdsp/dsp"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"syscall"
)

var log = slog.Scope("RadioServer")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
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
	//var frontend = frontends.CreateAirspyFrontend(0)
	var frontend = frontends.CreateLimeSDRFrontend(0)
	//var frontend = frontends.CreateTestSignalFrontend()
	frontend.Init()
	defer frontend.Destroy()

	frontend.SetCenterFrequency(106300000)
	//frontend.SetSampleRate(3000000)
	//frontend.SetGain(60)
	frontend.Start()

	defer frontend.Stop()

	log.Info("Protocol Version: %s", radioserver.ServerVersion.AsString())
	log.Info("SIMD Mode: %s", dsp.GetSIMDMode())

	srv := server.MakeRadioServer(frontend)
	err := srv.Listen(":4050", ":8000")
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
