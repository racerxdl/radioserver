package DSP

import (
	"github.com/quan-to/slog"
	"github.com/racerxdl/go.fifo"
	"github.com/racerxdl/radioserver/frontends"
	"github.com/racerxdl/radioserver/protocol"
	"github.com/racerxdl/radioserver/tools"
	"github.com/racerxdl/segdsp/dsp"
	"github.com/racerxdl/segdsp/dsp/fft"
	tools2 "github.com/racerxdl/segdsp/tools"
	"math"
	"runtime"
	"sync"
	"time"
)

var cgLog = slog.Scope("ChannelGenerator")

const maxFifoSize = 4096
const SmartFrameRate = 25
const SmartLength = 4096
const FFTAveraging = 2

type OnSmartIQSamples func(samples []complex64)
type OnIQSamples func(samples []complex64)
type OnFCSamples func(samples []float32)

type ChannelGenerator struct {
	sync.Mutex
	iqFrequencyTranslator    *dsp.FrequencyTranslator
	smartFrequencyTranslator *dsp.FrequencyTranslator
	fcFrequencyTranslator    *dsp.FrequencyTranslator

	inputFifo     *fifo.Queue
	running       bool
	settingsMutex sync.Mutex

	smartIQEnabled bool
	iqEnabled      bool
	fcEnabled      bool

	onIQSamples          OnIQSamples
	onSmartIQ            OnSmartIQSamples
	onFC                 OnFCSamples
	updateChannel        chan bool
	lastSmart            time.Time
	lastFC               time.Time
	smartIQPeriod        time.Duration
	blackmanWindow       []float32
	fcWindow             []float32
	lastFrequencySamples []float32

	syncSampleInput *sync.Cond

	smartIQSampleRate float32
	fcSampleRate      float32
	fcLength          uint32
}

func CreateChannelGenerator() *ChannelGenerator {
	var smartPeriod = 1e9 / float32(SmartFrameRate)

	var cg = &ChannelGenerator{
		Mutex:                sync.Mutex{},
		inputFifo:            fifo.NewQueue(),
		settingsMutex:        sync.Mutex{},
		updateChannel:        make(chan bool),
		lastSmart:            time.Now(),
		lastFC:               time.Now(),
		lastFrequencySamples: make([]float32, SmartLength),
		smartIQPeriod:        time.Duration(smartPeriod),
		fcLength:             SmartLength,
	}

	cg.syncSampleInput = sync.NewCond(cg)

	cg.blackmanWindow = make([]float32, SmartLength)
	w := dsp.BlackmanHarris(SmartLength, 92)
	for i, v := range w {
		cg.blackmanWindow[i] = float32(v)
	}

	cg.fcWindow = cg.blackmanWindow // Default size SmartIQ

	return cg
}

func (cg *ChannelGenerator) routine() {
	for cg.running {
		go func() {
			<-time.After(1 * time.Second)
			cg.syncSampleInput.Broadcast()
		}()
		cg.syncSampleInput.L.Lock()
		cg.syncSampleInput.Wait()
		cg.doWork()
		cg.syncSampleInput.L.Unlock()

		if !cg.running {
			break
		}
		runtime.Gosched()
	}
	cgLog.Debug("Cleaning fifo")
	for i := 0; i < cg.inputFifo.Len(); i++ {
		cg.inputFifo.Next()
	}
	cgLog.Debug("Done")
}

func (cg *ChannelGenerator) doWork() {
	cg.settingsMutex.Lock()
	for cg.inputFifo.Len() > 0 {
		var samples = cg.inputFifo.Next().([]complex64)
		if cg.iqEnabled {
			cg.processIQ(samples)
		}

		if cg.smartIQEnabled {
			cg.processSmart(samples)
		}

		if cg.fcEnabled {
			cg.processFrequency(samples)
		}
	}
	cg.settingsMutex.Unlock()
}

func (cg *ChannelGenerator) processIQ(samples []complex64) {
	if cg.onIQSamples != nil {
		if cg.iqFrequencyTranslator.GetDecimation() != 1 || cg.iqFrequencyTranslator.GetFrequency() != 0 {
			samples = cg.iqFrequencyTranslator.Work(samples)
		}
		cg.onIQSamples(samples)
	}
}

func (cg *ChannelGenerator) processFrequency(samples []complex64) {
	if time.Since(cg.lastFC) > cg.smartIQPeriod && cg.onFC != nil && len(samples) >= int(cg.fcLength) {
		// Process IQ Input
		if cg.fcFrequencyTranslator.GetDecimation() != 1 || cg.fcFrequencyTranslator.GetFrequency() != 0 {
			samples = cg.fcFrequencyTranslator.Work(samples)
		}

		samples = samples[:cg.fcLength]

		// Apply window to samples
		for j := 0; j < len(samples); j++ {
			var s = samples[j]
			var r = real(s) * float32(cg.fcWindow[j])
			var i = imag(s) * float32(cg.fcWindow[j])
			samples[j] = complex(r, i)
		}

		fftCData := fft.FFT(samples)

		var fftSamples = make([]float32, len(fftCData))
		var l = len(fftSamples)
		var lastV = float32(0)
		for i, v := range fftCData {
			var oI = (i + l/2) % l
			var m = float64(tools2.ComplexAbsSquared(v) * (1.0 / cg.fcSampleRate))

			m = 10 * math.Log10(m)

			fftSamples[oI] = (cg.lastFrequencySamples[i]*(FFTAveraging-1) + float32(m)) / FFTAveraging
			if fftSamples[i] != fftSamples[i] { // IsNaN
				fftSamples[i] = 0
			}

			if i > 0 {
				fftSamples[oI] = lastV*0.4 + fftSamples[oI]*0.6
			}

			lastV = fftSamples[oI]
		}

		copy(cg.lastFrequencySamples, fftSamples)

		cg.onFC(fftSamples)
		cg.lastFC = time.Now()
	}
}

func (cg *ChannelGenerator) processSmart(samples []complex64) {
	if time.Since(cg.lastSmart) > cg.smartIQPeriod && cg.onSmartIQ != nil {
		// Process IQ Input
		if cg.smartFrequencyTranslator.GetDecimation() != 1 || cg.smartFrequencyTranslator.GetFrequency() != 0 {
			samples = cg.smartFrequencyTranslator.Work(samples)
		}

		samples = samples[:SmartLength]

		// Apply window to samples
		for j := 0; j < len(samples); j++ {
			var s = samples[j]
			var r = real(s) * float32(cg.blackmanWindow[j])
			var i = imag(s) * float32(cg.blackmanWindow[j])
			samples[j] = complex(r, i)
		}

		cg.onSmartIQ(samples)
		cg.lastSmart = time.Now()
	}
}

func (cg *ChannelGenerator) notify() {
	cg.syncSampleInput.Broadcast()
}

func (cg *ChannelGenerator) Start() {
	if !cg.running {
		cgLog.Info("Starting Channel Generator")
		cg.running = true
		go cg.routine()
		//go func() {
		//	for cg.running {
		//		<-time.After(1 * time.Second)
		//		cgLog.Debug("Fifo Usage: %d", cg.inputFifo.UnsafeLen())
		//	}
		//}()
	}
}

func (cg *ChannelGenerator) Stop() {
	if cg.running {
		cgLog.Info("Stopping")
		cg.running = false
		cg.notify()
	}
}

func (cg *ChannelGenerator) StartIQ() {
	cg.settingsMutex.Lock()
	cgLog.Info("Enabling IQ")
	cg.iqEnabled = true
	cg.settingsMutex.Unlock()
}

func (cg *ChannelGenerator) StopIQ() {
	cg.settingsMutex.Lock()
	cgLog.Info("Disabling IQ")
	cg.iqEnabled = false
	cg.settingsMutex.Unlock()

	if !cg.smartIQEnabled && cg.running {
		go cg.Stop()
	}
}

func (cg *ChannelGenerator) StartFC() {
	cg.settingsMutex.Lock()
	cgLog.Info("Enabling Frequency Channel")
	cg.fcEnabled = true
	cg.settingsMutex.Unlock()
}

func (cg *ChannelGenerator) StopFC() {
	cg.settingsMutex.Lock()
	cgLog.Info("Disabling Frequency Channel")
	cg.fcEnabled = false
	cg.settingsMutex.Unlock()

	if !cg.smartIQEnabled && cg.running {
		go cg.Stop()
	}
}

func (cg *ChannelGenerator) StartSmartIQ() {
	cg.settingsMutex.Lock()
	cgLog.Info("Enabling SmartIQ")
	cg.smartIQEnabled = true
	cg.settingsMutex.Unlock()
}

func (cg *ChannelGenerator) StopSmartIQ() {
	cg.settingsMutex.Lock()
	cgLog.Info("Disabling SmartIQ")
	cg.smartIQEnabled = false

	if !cg.iqEnabled && cg.running {
		go cg.Stop()
	}
	cg.settingsMutex.Unlock()
}

func (cg *ChannelGenerator) UpdateSettings(channelType protocol.ChannelType, frontend frontends.Frontend, state *protocol.ChannelConfig) {
	cg.settingsMutex.Lock()
	cgLog.Info("Updating settings")

	var deviceFrequency = frontend.GetCenterFrequency()
	var deviceSampleRate = frontend.GetSampleRate()

	if channelType == protocol.ChannelType_IQ {
		var iqDecimationNumber = tools.StageToNumber(state.DecimationStage)
		var iqFtTaps = tools.GenerateTranslatorTaps(iqDecimationNumber, deviceSampleRate)
		var iqDeltaFrequency = float32(state.CenterFrequency) - float32(deviceFrequency)
		cgLog.Debug("IQ Delta Frequency: %.0f", iqDeltaFrequency)
		cg.iqFrequencyTranslator = dsp.MakeFrequencyTranslator(int(iqDecimationNumber), iqDeltaFrequency, float32(deviceSampleRate), iqFtTaps)
	}

	if channelType == protocol.ChannelType_SmartIQ {
		var smartIQDecimationNumber = tools.StageToNumber(state.DecimationStage)
		var smartFtTaps = tools.GenerateTranslatorTaps(smartIQDecimationNumber, deviceSampleRate)
		var smartIQDeltaFrequency = float32(state.CenterFrequency) - float32(deviceFrequency)
		cgLog.Debug("SmartIQ Delta Frequency: %.0f", smartIQDeltaFrequency)
		cg.smartFrequencyTranslator = dsp.MakeFrequencyTranslator(int(smartIQDecimationNumber), smartIQDeltaFrequency, float32(deviceSampleRate), smartFtTaps)
		cg.smartIQSampleRate = float32(deviceSampleRate / smartIQDecimationNumber)
	}

	cg.settingsMutex.Unlock()
	cgLog.Info("Settings updated.")
}

func (cg *ChannelGenerator) UpdateFrequencyChannel(frontend frontends.Frontend, state *protocol.FrequencyChannelConfig) {
	cg.settingsMutex.Lock()
	cgLog.Info("Updating Frequency Channel Settings")

	var deviceFrequency = frontend.GetCenterFrequency()
	var deviceSampleRate = frontend.GetSampleRate()

	var fcDecimationNumber = tools.StageToNumber(state.DecimationStage)
	var fcFtTaps = tools.GenerateTranslatorTaps(fcDecimationNumber, deviceSampleRate)
	var fcDeltaFrequency = float32(state.CenterFrequency) - float32(deviceFrequency)
	cgLog.Debug("FC Delta Frequency: %.0f", fcDeltaFrequency)
	cg.fcFrequencyTranslator = dsp.MakeFrequencyTranslator(int(fcDecimationNumber), fcDeltaFrequency, float32(deviceSampleRate), fcFtTaps)
	cg.fcSampleRate = float32(deviceSampleRate / fcDecimationNumber)

	if cg.fcLength != state.Length {
		cg.fcLength = state.Length
		cg.fcWindow = make([]float32, cg.fcLength)
		w := dsp.BlackmanHarris(int(cg.fcLength), 92)
		for i, v := range w {
			cg.fcWindow[i] = float32(v)
		}
		cg.lastFrequencySamples = make([]float32, cg.fcLength)
	}

	cg.settingsMutex.Unlock()
	cgLog.Info("Settings updated.")
}

func (cg *ChannelGenerator) PushSamples(samples []complex64) {
	if !cg.running {
		return
	}

	var fifoLength = cg.inputFifo.Len()

	if maxFifoSize <= fifoLength {
		cgLog.Debug("Fifo Overflowing!")
		return
	}

	cg.inputFifo.Add(samples)

	cg.notify()
}

func (cg *ChannelGenerator) SetOnIQ(cb OnIQSamples) {
	cg.onIQSamples = cb
}

func (cg *ChannelGenerator) SetOnSmartIQ(cb OnSmartIQSamples) {
	cg.onSmartIQ = cb
}

func (cg *ChannelGenerator) SetOnFC(cb OnFCSamples) {
	cg.onFC = cb
}

func (cg *ChannelGenerator) SmartIQRunning() bool {
	return cg.smartIQEnabled
}

func (cg *ChannelGenerator) IQRunning() bool {
	return cg.iqEnabled
}

func (cg *ChannelGenerator) GetSmartIQSampleRate() float32 {
	cg.settingsMutex.Lock()
	sr := cg.smartIQSampleRate
	cg.settingsMutex.Unlock()

	return sr
}
