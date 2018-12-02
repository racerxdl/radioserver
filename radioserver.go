package main

import (
	"flag"
	"fmt"
	"github.com/racerxdl/radioserver/SLog"
	"github.com/racerxdl/radioserver/frontends"
	"github.com/racerxdl/radioserver/protocol"
	"github.com/racerxdl/segdsp/dsp"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
	"syscall"
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			SLog.Fatal(err)
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

	hash, _ := strconv.ParseInt(commitHash, 16, 32)
	ServerVersion.Hash = uint32(hash)

	SLog.Info("Protocol Version: %s", ServerVersion.String())
	SLog.Info("Commit Hash: %s", commitHash)
	SLog.Info("SIMD Mode: %s", dsp.GetSIMDMode())

	var frontend = frontends.CreateAirspyFrontend(0)
	//var frontend = frontends.CreateLimeSDRFrontend(0)
	frontend.Init()
	frontend.SetCenterFrequency(97700000)

	defer frontend.Destroy()

	var name = frontend.GetShortName()

	SLog.Info("Frontend: %s", frontend.GetName())

	var deviceName [16]uint8

	for i := 0; i < 15; i++ {
		deviceName[i] = name[i]
	}
	deviceName[15] = 0x00

	serverState.Frontend = frontend
	serverState.CanControl = 0

	serverState.DeviceInfo = protocol.DeviceInfo{
		DeviceType:        frontend.GetDeviceType(),
		DeviceSerial:      frontend.GetUintDeviceSerial(),
		MaximumSampleRate: frontend.GetMaximumSampleRate(),
		DecimationStages:  frontend.MaximumDecimationStages(),
		MaximumGainValue:  frontend.MaximumGainValue(),
		MinimumFrequency:  frontend.MinimumFrequency(),
		MaximumFrequency:  frontend.MaximumFrequency(),
		Resolution:        uint32(frontend.GetResolution()),
		DeviceName:        deviceName,
	}

	frontend.SetSamplesAvailableCallback(serverState.PushSamples)

	stop := make(chan bool, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		SLog.Info("Got SIGTERM! Closing it")
		tcpServerStatus = false
		stop <- true
	}()

	runServer(stop)
	SLog.Info("Closing")
}
