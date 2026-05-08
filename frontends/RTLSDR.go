//go:build rtlsdr

package frontends

import (
	"fmt"
	"github.com/jpoirier/gortlsdr"
	"github.com/quan-to/slog"
	"github.com/racerxdl/radioserver/protocol"
	"math"
)

const rtlsdrMinimumFrequency uint32 = 24e6
const rtlsdrMaximumFrequency uint32 = 1766e6
const rtlsdrMinSampleRate uint32 = 10e3

var rtlsdrLog = slog.Scope("RTLSDR Frontend")

type RTLSDRFrontend struct {
	device             *rtlsdr.Context
	cb                 SamplesCallback
	deviceIndex        int
	deviceSerial       string
	maxSampleRate      uint32
	maxDecimationStage uint32
	currentGain        uint8
	running            bool
}

func CreateRTLSDRFrontend(index int) Frontend {
	return &RTLSDRFrontend{
		deviceIndex:  index,
		currentGain:  0,
		running:      false,
		deviceSerial: "",
	}
}

func (f *RTLSDRFrontend) Init() bool {
	dev, err := rtlsdr.Open(f.deviceIndex)
	if err != nil {
		rtlsdrLog.Error("Failed to open device %d: %s", f.deviceIndex, err)
		return false
	}
	f.device = dev

	_, product, serial, _ := dev.GetUsbStrings()
	f.deviceSerial = fmt.Sprintf("%s:%s:%s", rtlsdr.GetDeviceName(f.deviceIndex), product, serial)

	f.maxSampleRate = 2400000
	f.maxDecimationStage = f.calcMaxDecimation()

	rtlsdrLog.Info("Opened RTLSDR: %s (serial: %s)", rtlsdr.GetDeviceName(f.deviceIndex), serial)
	return true
}

func (f *RTLSDRFrontend) calcMaxDecimation() uint32 {
	var stage uint32
	calcSR := f.maxSampleRate
	for calcSR >= rtlsdrMinSampleRate {
		stage++
		calcSR = f.maxSampleRate / uint32(math.Pow(2, float64(stage)))
	}
	return stage
}

func (f *RTLSDRFrontend) GetDeviceType() protocol.DeviceType {
	return protocol.DeviceRtlsdr
}

func (f *RTLSDRFrontend) GetDeviceSerial() string {
	return f.deviceSerial
}

func (f *RTLSDRFrontend) GetUintDeviceSerial() uint32 {
	return 0
}

func (f *RTLSDRFrontend) GetName() string {
	return rtlsdr.GetDeviceName(f.deviceIndex)
}

func (f *RTLSDRFrontend) GetShortName() string {
	return "RTLSDR"
}

func (f *RTLSDRFrontend) GetMaximumSampleRate() uint32 {
	return f.maxSampleRate
}

func (f *RTLSDRFrontend) GetSampleRate() uint32 {
	if f.device == nil {
		return 0
	}
	return uint32(f.device.GetSampleRate())
}

func (f *RTLSDRFrontend) SetSampleRate(sampleRate uint32) uint32 {
	if f.device == nil {
		return 0
	}
	_ = f.device.SetSampleRate(int(sampleRate))
	return uint32(f.device.GetSampleRate())
}

func (f *RTLSDRFrontend) GetAvailableSampleRates() []uint32 {
	return []uint32{2400000, 2048000, 1024000, 320000}
}

func (f *RTLSDRFrontend) GetCenterFrequency() uint32 {
	if f.device == nil {
		return 0
	}
	return uint32(f.device.GetCenterFreq())
}

func (f *RTLSDRFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	if f.device == nil {
		return 0
	}
	_ = f.device.SetCenterFreq(int(centerFrequency))
	return uint32(f.device.GetCenterFreq())
}

func (f *RTLSDRFrontend) MinimumFrequency() uint32 {
	return rtlsdrMinimumFrequency
}

func (f *RTLSDRFrontend) MaximumFrequency() uint32 {
	return rtlsdrMaximumFrequency
}

func (f *RTLSDRFrontend) GetMaximumBandwidth() uint32 {
	return uint32(float32(f.maxSampleRate) * 0.8)
}

func (f *RTLSDRFrontend) MaximumGainValue() uint32 {
	if f.device == nil {
		return 0
	}
	gains, err := f.device.GetTunerGains()
	if err != nil || len(gains) == 0 {
		return 0
	}
	return uint32(gains[len(gains)-1] / 10)
}

func (f *RTLSDRFrontend) MaximumDecimationStages() uint32 {
	return f.maxDecimationStage
}

func (f *RTLSDRFrontend) GetGain() uint8 {
	return f.currentGain
}

func (f *RTLSDRFrontend) SetGain(value uint8) {
	if f.device == nil {
		return
	}
	_ = f.device.SetTunerGainMode(true)
	_ = f.device.SetTunerGain(int(value) * 10)
	f.currentGain = value
}

func (f *RTLSDRFrontend) SetAGC(agc bool) {
	if f.device == nil {
		return
	}
	_ = f.device.SetTunerGainMode(!agc)
	_ = f.device.SetAgcMode(agc)
}

func (f *RTLSDRFrontend) SetAntenna(value string) {
	rtlsdrLog.Warn("RTLSDR does not support antenna switch. Ignoring...")
}

func (f *RTLSDRFrontend) SetBiasT(value bool) {
	if f.device == nil {
		return
	}
	_ = f.device.SetBiasTee(value)
}

func (f *RTLSDRFrontend) GetResolution() uint8 {
	return 8
}

func (f *RTLSDRFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.cb = cb
}

func (f *RTLSDRFrontend) Start() {
	if f.running || f.device == nil {
		return
	}
	rtlsdrLog.Info("Starting")
	_ = f.device.ResetBuffer()
	f.running = true

	go func() {
		cb := func(buf []byte) {
			if !f.running || f.cb == nil {
				return
			}
			n := len(buf) / 2
			samples := make([]complex64, n)
			for i := 0; i < n; i++ {
				re := (float32(buf[i*2]) - 127.5) / 127.5
				im := (float32(buf[i*2+1]) - 127.5) / 127.5
				samples[i] = complex(re, im)
			}
			f.cb(samples)
		}

		err := f.device.ReadAsync(cb, nil, 0, 0)
		if err != nil {
			rtlsdrLog.Error("ReadAsync error: %s", err)
		}
	}()
}

func (f *RTLSDRFrontend) Stop() {
	if !f.running || f.device == nil {
		return
	}
	rtlsdrLog.Info("Stopping")
	f.running = false
	_ = f.device.CancelAsync()
}

func (f *RTLSDRFrontend) Destroy() {
	if f.device == nil {
		return
	}
	rtlsdrLog.Info("De-initializing")
	f.Stop()
	_ = f.device.Close()
	f.device = nil
}
