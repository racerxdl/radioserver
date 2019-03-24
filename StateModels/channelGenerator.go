package StateModels

import (
	"github.com/quan-to/slog"
	"github.com/racerxdl/go.fifo"
	"github.com/racerxdl/radioserver/protocol"
	"github.com/racerxdl/radioserver/tools"
	"github.com/racerxdl/segdsp/dsp"
	"sync"
	"time"
)

var cgLog = slog.Scope("ChannelGenerator")

const maxFifoSize = 4096
const SmartFrameRate = 20
const SmartLength = 4096

type OnSmartIQSamples func(samples []complex64)
type OnIQSamples func(samples []complex64)

type ChannelGenerator struct {
	sync.Mutex
	iqFrequencyTranslator    *dsp.FrequencyTranslator
	smartFrequencyTranslator *dsp.FrequencyTranslator

	inputFifo     *fifo.Queue
	running       bool
	settingsMutex sync.Mutex

	smartIQEnabled bool
	iqEnabled      bool

	onIQSamples    OnIQSamples
	onSmartIQ      OnSmartIQSamples
	updateChannel  chan bool
	lastSmart      time.Time
	smartIQPeriod  time.Duration
	blackmanWindow []float32

	syncSampleInput *sync.Cond
}

func CreateChannelGenerator() *ChannelGenerator {
	var smartPeriod = 1e9 / float32(SmartFrameRate)

	var cg = &ChannelGenerator{
		Mutex:         sync.Mutex{},
		inputFifo:     fifo.NewQueue(),
		settingsMutex: sync.Mutex{},
		updateChannel: make(chan bool),
		lastSmart:     time.Now(),
		smartIQPeriod: time.Duration(smartPeriod),
	}

	cg.syncSampleInput = sync.NewCond(cg)

	cg.blackmanWindow = make([]float32, SmartLength)
	w := dsp.BlackmanHarris(SmartLength, 92)
	for i, v := range w {
		cg.blackmanWindow[i] = float32(v)
	}

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
	}
	cgLog.Debug("Cleaning fifo")
	for i := 0; i < cg.inputFifo.UnsafeLen(); i++ {
		cg.inputFifo.UnsafeNext()
	}
	cgLog.Debug("Done")
}

func (cg *ChannelGenerator) doWork() {
	cg.settingsMutex.Lock()
	defer cg.settingsMutex.Unlock()

	for cg.inputFifo.UnsafeLen() > 0 {
		var samples = cg.inputFifo.Next().([]complex64)
		if cg.iqEnabled {
			cg.processIQ(samples)
		}

		if cg.smartIQEnabled {
			cg.processSmart(samples)
		}
	}
}

func (cg *ChannelGenerator) processIQ(samples []complex64) {
	if cg.onIQSamples != nil {
		if cg.iqFrequencyTranslator.GetDecimation() != 1 || cg.iqFrequencyTranslator.GetFrequency() != 0 {
			samples = cg.iqFrequencyTranslator.Work(samples)
		}
		cg.onIQSamples(samples)
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
		if cg.iqFrequencyTranslator == nil && cg.smartFrequencyTranslator == nil {
			cgLog.Fatal("Trying to start Channel Generator without frequencyTranslator for either IQ or Smart")
		}
		cg.running = true
		go cg.routine()
		go func() {
			for cg.running {
				<-time.After(1 * time.Second)
				cgLog.Debug("Fifo Usage: %d", cg.inputFifo.UnsafeLen())
			}
		}()
	}
}

func (cg *ChannelGenerator) Stop() {
	if cg.running {
		cgLog.Info("Stopping")
		cg.running = false
		cg.notify()
	}
}

func (cg *ChannelGenerator) UpdateSettings(state *ClientState) {
	cg.settingsMutex.Lock()
	cgLog.Info("Updating settings")

	var deviceFrequency = state.ServerState.Frontend.GetCenterFrequency()
	var deviceSampleRate = state.ServerState.Frontend.GetSampleRate()

	cg.iqEnabled = state.CGS.StreamingMode == protocol.TypeIQ || state.CGS.StreamingMode == protocol.TypeCombined
	cg.smartIQEnabled = state.CGS.StreamingMode == protocol.TypeSmartIQ || state.CGS.StreamingMode == protocol.TypeCombined

	// region IQ Channel
	if cg.iqEnabled {
		var iqDecimationNumber = tools.StageToNumber(state.CGS.IQDecimation)
		var iqFtTaps = tools.GenerateTranslatorTaps(iqDecimationNumber, deviceSampleRate)
		var iqDeltaFrequency = float32(state.CGS.IQCenterFrequency) - float32(deviceFrequency)
		cgLog.Debug("IQ Delta Frequency: %.0f", iqDeltaFrequency)
		cg.iqFrequencyTranslator = dsp.MakeFrequencyTranslator(int(iqDecimationNumber), iqDeltaFrequency, float32(deviceSampleRate), iqFtTaps)
	}
	// endregion
	// region Smart IQ Channel
	if cg.smartIQEnabled {
		var smartIQDecimationNumber = tools.StageToNumber(state.CGS.SmartIQDecimation)
		var smartFtTaps = tools.GenerateTranslatorTaps(smartIQDecimationNumber, deviceSampleRate)
		var smartIQDeltaFrequency = float32(state.CGS.SmartCenterFrequency) - float32(deviceFrequency)
		cgLog.Debug("SmartIQ Delta Frequency: %.0f", smartIQDeltaFrequency)
		cg.smartFrequencyTranslator = dsp.MakeFrequencyTranslator(int(smartIQDecimationNumber), smartIQDeltaFrequency, float32(deviceSampleRate), smartFtTaps)
	}
	// endregion
	cg.settingsMutex.Unlock()
	if state.CGS.Streaming && !cg.running {
		cg.Start()
	}

	if !state.CGS.Streaming && cg.running {
		cg.Stop()
	}
	cgLog.Info("Settings updated.")
}

func (cg *ChannelGenerator) PushSamples(samples []complex64) {
	if !cg.running {
		return
	}

	var fifoLength = cg.inputFifo.UnsafeLen()

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
