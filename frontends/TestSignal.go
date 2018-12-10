package frontends

import (
	"fmt"
	"github.com/racerxdl/radioserver/SLog"
	"github.com/racerxdl/radioserver/protocol"
	"github.com/racerxdl/segdsp/dsp"
	"math"
	"math/rand"
	"time"
)

var testSignalSampleRate = 10e6
var testSignalLog = SLog.Scope("TestSignal Frontend")

type TestSignalFrontend struct {
	cb SamplesCallback

	deviceSerial uint64
	currentGain  uint8
	running      bool

	availableSampleRates []uint32
	samplesBuffer        []complex64
	frequency            uint32
}

func CreateTestSignalFrontend() Frontend {
	var f = &TestSignalFrontend{
		deviceSerial: 0,
		currentGain:  0,
		running:      false,
		frequency:    0,
	}

	// region Cache Samples
	var samples = make([]complex64, 1024)
	var interp = dsp.MakeInterpolator(16)
	var lowPass = dsp.MakeLowPassFixed(1, testSignalSampleRate, testSignalSampleRate/2-5e3, 63)
	var frequencyShift = dsp.MakeFrequencyTranslator(1, -100e3, float32(testSignalSampleRate), lowPass)

	for i := 0; i < len(samples); i++ {
		samples[i] = complex((rand.Float32()-1)*0.5, 0)
	}

	samples = interp.Work(samples)

	// Generate some background noise
	for i := 0; i < len(samples); i++ {
		samples[i] += complex((rand.Float32()-1)*1e-4, 0)
	}

	samples = frequencyShift.Work(samples)
	f.samplesBuffer = make([]complex64, 16384)
	copy(f.samplesBuffer, samples)
	// endregion

	return f
}

func (f *TestSignalFrontend) GetUintDeviceSerial() uint32 {
	return uint32(f.deviceSerial & 0xFFFFFFFF)
}

func (f *TestSignalFrontend) MinimumFrequency() uint32 {
	return limeMinimumFrequency
}

func (f *TestSignalFrontend) MaximumFrequency() uint32 {
	return limeMaximumFrequency
}

func (f *TestSignalFrontend) GetMaximumBandwidth() uint32 {
	return uint32(testSignalSampleRate)
}

func (f *TestSignalFrontend) MaximumGainValue() uint32 {
	return 32
}

func (f *TestSignalFrontend) MaximumDecimationStages() uint32 {
	return 8
}

func (f *TestSignalFrontend) GetDeviceType() uint32 {
	return protocol.DeviceTestSignal
}

func (f *TestSignalFrontend) GetDeviceSerial() string {
	return fmt.Sprintf("%08x", rand.Uint32())
}
func (f *TestSignalFrontend) GetMaximumSampleRate() uint32 {
	return uint32(testSignalSampleRate)
}
func (f *TestSignalFrontend) SetSampleRate(sampleRate uint32) uint32 {
	return uint32(testSignalSampleRate)
}
func (f *TestSignalFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	f.frequency = centerFrequency
	return f.frequency
}
func (f *TestSignalFrontend) GetAvailableSampleRates() []uint32 {
	return f.availableSampleRates
}

func (f *TestSignalFrontend) loop() {
	interval := time.Duration((1e9 * float64(len(f.samplesBuffer))) / testSignalSampleRate)
	loopTicker := time.NewTicker(interval)
	testSignalLog.Debug("Period: %v", interval)
	for f.running {
		for range loopTicker.C {
			f.work()
		}
	}
}

func (f *TestSignalFrontend) work() {
	if f.cb != nil {
		var samples = make([]complex64, len(f.samplesBuffer))
		copy(samples, f.samplesBuffer)
		if f.currentGain > 0 {
			var aGain = float32(math.Pow(10, float64(f.currentGain)/10))
			for j := 0; j < len(samples); j++ {
				var r = real(samples[j])
				var i = imag(samples[j])
				samples[j] = complex(r*aGain, i*aGain)
			}
		}
		f.cb(samples)
	}
}

func (f *TestSignalFrontend) Start() {
	if !f.running {
		testSignalLog.Info("Starting")
		f.running = true
		go f.loop()
	}
}
func (f *TestSignalFrontend) Stop() {
	if f.running {
		testSignalLog.Info("Stopping")
		f.running = false
	}
}
func (f *TestSignalFrontend) SetAntenna(value string) {}
func (f *TestSignalFrontend) SetAGC(agc bool)         {}
func (f *TestSignalFrontend) SetGain(value uint8) {
	f.currentGain = value
}
func (f *TestSignalFrontend) GetGain() uint8 {
	return f.currentGain
}
func (f *TestSignalFrontend) SetBiasT(value bool) {}
func (f *TestSignalFrontend) GetCenterFrequency() uint32 {
	return uint32(f.frequency)
}
func (f *TestSignalFrontend) GetName() string {
	return "Test Signal Generator"
}
func (f *TestSignalFrontend) GetShortName() string {
	return "TestSignal"
}
func (f *TestSignalFrontend) GetSampleRate() uint32 {
	return uint32(testSignalSampleRate)
}
func (f *TestSignalFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.cb = cb
}

func (f *TestSignalFrontend) GetResolution() uint8 {
	return 32
}

func (f *TestSignalFrontend) Init() bool {
	return true
}

func (f *TestSignalFrontend) Destroy() {}
