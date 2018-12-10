package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/racerxdl/radioserver/client"
	"github.com/racerxdl/radioserver/protocol"
	"log"
	"os"
	"time"
)

var f *os.File

type MyCallback struct{}

func (cb *MyCallback) OnData(dType int, data interface{}) {
	if dType == client.SamplesComplex64 {
		samples := data.([]complex64)
		log.Println("Received Complex 64 bit Data! ", len(samples))
	} else if dType == client.SamplesComplex32 {
		samples := data.([]client.ComplexInt16)
		log.Println("Received Complex 32 bit Data! ", len(samples))
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, data)

		_, _ = f.Write(buf.Bytes())
	} else if dType == client.SmartSamplesComplex32 {
		log.Println("Got Smart IQ 32 bit Samples!")
	} else if dType == client.SmartSamplesComplex64 {
		log.Println("Got Smart IQ 64 bit Samples!")
	} else if dType == client.DeviceSync {
		log.Println("Got device sync!")
	}
}

func main() {
	var rs = client.MakeRadioClient("127.0.0.1", protocol.DefaultPort)

	var cb = MyCallback{}

	f, _ = os.Create("iq.raw")

	defer f.Close()

	rs.SetCallback(&cb)

	rs.Connect()

	log.Println(fmt.Sprintf("Device: %s", rs.GetName()))
	var srs = rs.GetAvailableSampleRates()

	log.Println("Available SampleRates:")
	for i := 0; i < len(srs); i++ {
		log.Println(fmt.Sprintf("		%f msps", float32(srs[i])/1e6))
	}
	if rs.SetSampleRate(2500000) == protocol.Invalid {
		log.Println("Error setting sample rate.")
	}
	if rs.SetCenterFrequency(106300000) == protocol.Invalid {
		log.Println("Error setting center frequency.")
	}

	rs.SetStreamingMode(protocol.TypeIQ)

	log.Println("Starting")
	rs.Start()

	time.Sleep(time.Second * 10)

	log.Print("Stopping")
	rs.Stop()

	rs.Disconnect()
}
