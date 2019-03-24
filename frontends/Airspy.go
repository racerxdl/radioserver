package frontends

import (
	"fmt"
	"github.com/quan-to/slog"
	"github.com/racerxdl/radioserver/protocol"
	"github.com/racerxdl/spy2go/airspy"
	"github.com/racerxdl/spy2go/spytypes"
	"math"
)

const airspyMaximumFrequency = 1.768e6
const airspyMinimumFrequency = 24e6

var airspyLog = slog.Scope("Airspy Frontend")

type AirspyFrontend struct {
	device *airspy.Device
	cb     SamplesCallback

	deviceSerial       uint64
	maxSampleRate      uint32
	maxDecimationStage uint32
	currentGain        uint8
	running            bool
}

type internalCallback struct {
	parent func(dType int, data interface{})
}

func (ic *internalCallback) OnData(dType int, data interface{}) {
	ic.parent(dType, data)
}

func CreateAirspyFrontend(serial uint64) Frontend {
	airspy.Initialize()
	var f = &AirspyFrontend{
		device:        airspy.MakeAirspyDevice(serial),
		deviceSerial:  0,
		maxSampleRate: 0,
		currentGain:   0,
		running:       false,
	}

	f.device.SetSampleType(spytypes.SamplesComplex64)

	if f.deviceSerial == 0 {
		// Fetch device serial
		f.deviceSerial = f.device.GetSerial()
	}

	for _, v := range f.device.GetAvailableSampleRates() {
		if v > f.maxSampleRate {
			f.maxSampleRate = v
		}
	}

	var maxDecimationStage = uint32(0)
	var calcSR = f.maxSampleRate

	for calcSR >= minimumSampleRate {
		maxDecimationStage += 1
		var decim = uint32(math.Pow(2, float64(maxDecimationStage)))
		calcSR = f.maxSampleRate / decim
	}

	f.maxDecimationStage = maxDecimationStage

	var ic = &internalCallback{
		parent: f.internalCb,
	}

	f.device.SetCallback(ic)
	f.device.SetSampleRate(f.maxSampleRate)

	return f
}

func (f *AirspyFrontend) GetUintDeviceSerial() uint32 {
	return uint32(f.deviceSerial & 0xFFFFFFFF)
}

func (f *AirspyFrontend) MinimumFrequency() uint32 {
	return airspyMinimumFrequency
}

func (f *AirspyFrontend) MaximumFrequency() uint32 {
	return airspyMaximumFrequency
}

func (f *AirspyFrontend) GetMaximumBandwidth() uint32 {
	return uint32(float32(f.maxSampleRate) * 0.8)
}

func (f *AirspyFrontend) MaximumGainValue() uint32 {
	return 16
}

func (f *AirspyFrontend) MaximumDecimationStages() uint32 {
	return f.maxDecimationStage
}

func (f *AirspyFrontend) GetDeviceType() uint32 {
	return protocol.DeviceAirspyOne
}

func (f *AirspyFrontend) internalCb(dType int, data interface{}) {
	if dType != spytypes.SamplesComplex64 {
		panic("Spy2Go Library is sending different types than we asked!")
	}

	if f.cb != nil {
		f.cb(data.([]complex64))
	}
}

func (f *AirspyFrontend) GetDeviceSerial() string {
	return fmt.Sprintf("%08x", f.deviceSerial)
}
func (f *AirspyFrontend) GetMaximumSampleRate() uint32 {
	return f.maxSampleRate
}
func (f *AirspyFrontend) SetSampleRate(sampleRate uint32) uint32 {
	f.device.SetSampleRate(sampleRate)
	return f.device.GetSampleRate()
}
func (f *AirspyFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	f.device.SetCenterFrequency(centerFrequency)
	return f.device.GetCenterFrequency()
}
func (f *AirspyFrontend) GetAvailableSampleRates() []uint32 {
	return f.device.GetAvailableSampleRates()
}
func (f *AirspyFrontend) Start() {
	if !f.running {
		airspyLog.Info("Starting")
		f.device.Start()
		f.running = true
	}
}
func (f *AirspyFrontend) Stop() {
	if f.running {
		airspyLog.Info("Stopping")
		f.device.Stop()
		f.running = false
	}
}
func (f *AirspyFrontend) SetAntenna(value string) {
	airspyLog.Warn("Airspy Frontend does not support antenna switch. Ignoring...")
}
func (f *AirspyFrontend) SetAGC(agc bool) {
	f.device.SetAGC(agc)
}
func (f *AirspyFrontend) SetGain(value uint8) {
	f.device.SetLinearityGain(value)
	f.currentGain = value
}
func (f *AirspyFrontend) GetGain() uint8 {
	return f.currentGain
}
func (f *AirspyFrontend) SetBiasT(value bool) {
	f.device.SetBiasT(value)
}
func (f *AirspyFrontend) GetCenterFrequency() uint32 {
	return f.device.GetCenterFrequency()
}
func (f *AirspyFrontend) GetName() string {
	return f.device.GetName()
}
func (f *AirspyFrontend) GetShortName() string {
	return "Airspy"
}
func (f *AirspyFrontend) GetSampleRate() uint32 {
	return f.device.GetSampleRate()
}
func (f *AirspyFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.cb = cb
}
func (f *AirspyFrontend) Init() bool {
	return true
}

func (f *AirspyFrontend) GetResolution() uint8 {
	return 12
}

func (f *AirspyFrontend) Destroy() {
	airspyLog.Info("De-initializing")
	airspy.DeInitialize()
}
