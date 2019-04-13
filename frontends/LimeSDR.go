package frontends

import (
	"fmt"
	"github.com/myriadrf/limedrv"
	"github.com/quan-to/slog"
	"github.com/racerxdl/radioserver/protocol"
	"math"
)

const limeMaximumFrequency = 3.8e6
const limeMinimumFrequency = 100e3

var limeLog = slog.Scope("LimeSDR Frontend")

type LimeSDRFrontend struct {
	device *limedrv.LMSDevice
	cb     SamplesCallback

	deviceSerial       uint64
	maxSampleRate      uint32
	maxDecimationStage uint32
	currentGain        uint8
	running            bool

	availableSampleRates []uint32
	selectedChannel      *limedrv.LMSChannel
	selectedChannelIndex int
	//selectedAntenna      int
}

func CreateLimeSDRFrontend(deviceIdx int) Frontend {
	devices := limedrv.GetDevices()
	if len(devices) == 0 {
		limeLog.Fatal("No devices found.\n")
	}

	if len(devices) < deviceIdx {
		limeLog.Fatal("No such device %d found.\n", deviceIdx)
	}

	var device = limedrv.Open(devices[deviceIdx])

	var f = &LimeSDRFrontend{
		device:               device,
		deviceSerial:         0,
		maxSampleRate:        0,
		currentGain:          0,
		running:              false,
		selectedChannelIndex: 0,
	}

	f.deviceSerial = 0
	f.maxSampleRate = 5000000 //60000000

	var availableSampleRates = make([]uint32, 1)
	availableSampleRates[0] = f.maxSampleRate

	var maxDecimationStage = uint32(0)
	var calcSR = f.maxSampleRate

	for calcSR >= minimumSampleRate {
		maxDecimationStage += 1
		var decim = uint32(math.Pow(2, float64(maxDecimationStage)))
		calcSR = f.maxSampleRate / decim
		availableSampleRates = append(availableSampleRates, calcSR)
	}

	f.availableSampleRates = availableSampleRates

	f.maxDecimationStage = maxDecimationStage

	f.device.
		SetCallback(func(samples []complex64, _ int, _ uint64) {
			if f.cb != nil {
				f.cb(samples)
			}
		})
	f.device.SetSampleRate(float64(f.maxSampleRate), 2)

	f.selectedChannel = device.RXChannels[f.selectedChannelIndex]

	f.selectedChannel.
		Enable().
		SetLPF(float64(f.maxSampleRate) / 2).
		EnableLPF().
		SetDigitalLPF(float64(f.maxSampleRate) / 2).
		EnableDigitalLPF().
		SetAntennaByName("LNAW")

	f.device.SetGainNormalized(f.selectedChannelIndex, true, 0.4)

	return f
}

func (f *LimeSDRFrontend) GetUintDeviceSerial() uint32 {
	return uint32(f.deviceSerial & 0xFFFFFFFF)
}

func (f *LimeSDRFrontend) MinimumFrequency() uint32 {
	return limeMinimumFrequency
}

func (f *LimeSDRFrontend) MaximumFrequency() uint32 {
	return limeMaximumFrequency
}

func (f *LimeSDRFrontend) GetMaximumBandwidth() uint32 {
	return f.maxSampleRate
}

func (f *LimeSDRFrontend) MaximumGainValue() uint32 {
	return 32
}

func (f *LimeSDRFrontend) MaximumDecimationStages() uint32 {
	return f.maxDecimationStage
}

func (f *LimeSDRFrontend) GetDeviceType() protocol.DeviceType {
	return protocol.DeviceLimeSDRUSB
}

func (f *LimeSDRFrontend) GetDeviceSerial() string {
	return fmt.Sprintf("%08x", f.deviceSerial)
}
func (f *LimeSDRFrontend) GetMaximumSampleRate() uint32 {
	return f.maxSampleRate
}
func (f *LimeSDRFrontend) SetSampleRate(sampleRate uint32) uint32 {
	var overSample = 2 * (f.maxSampleRate / sampleRate)
	f.device.SetSampleRate(float64(sampleRate), int(overSample))
	deviceSr, _ := f.device.GetSampleRate()
	return uint32(deviceSr)
}
func (f *LimeSDRFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	f.device.SetCenterFrequency(f.selectedChannelIndex, true, float64(centerFrequency))
	return uint32(f.device.GetCenterFrequency(f.selectedChannelIndex, true))
}
func (f *LimeSDRFrontend) GetAvailableSampleRates() []uint32 {
	return f.availableSampleRates
}
func (f *LimeSDRFrontend) Start() {
	if !f.running {
		limeLog.Info("Starting")
		f.device.Start()
		f.running = true
	}
}
func (f *LimeSDRFrontend) Stop() {
	if f.running {
		limeLog.Info("Stopping")
		f.device.Stop()
		f.running = false
	}
}
func (f *LimeSDRFrontend) SetAntenna(value string) {
	//limeLog.Warn("Airspy Frontend does not support antenna switch. Ignoring...")
}
func (f *LimeSDRFrontend) SetAGC(agc bool) {
	//f.device.SetAGC(agc)
}
func (f *LimeSDRFrontend) SetGain(value uint8) {
	caculatedGain := float64(f.MaximumGainValue()) * (float64(value) / 256)
	f.device.SetGainNormalized(f.selectedChannelIndex, true, caculatedGain)
	f.currentGain = value
}
func (f *LimeSDRFrontend) GetGain() uint8 {
	return f.currentGain
}
func (f *LimeSDRFrontend) SetBiasT(value bool) {
	//f.device.SetBiasT(value)
}
func (f *LimeSDRFrontend) GetCenterFrequency() uint32 {
	return uint32(f.device.GetCenterFrequency(f.selectedChannelIndex, true))
}
func (f *LimeSDRFrontend) GetName() string {
	return "LimeSDR USB"
}
func (f *LimeSDRFrontend) GetShortName() string {
	return "LimeSDR"
}
func (f *LimeSDRFrontend) GetSampleRate() uint32 {
	deviceSr, _ := f.device.GetSampleRate()
	return uint32(deviceSr)
}
func (f *LimeSDRFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.cb = cb
}

func (f *LimeSDRFrontend) GetResolution() uint8 {
	return 12
}

func (f *LimeSDRFrontend) Init() bool {
	return true
}

func (f *LimeSDRFrontend) Destroy() {

}
