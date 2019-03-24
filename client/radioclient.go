package client

import (
	"context"
	"github.com/quan-to/slog"
	"github.com/racerxdl/radioserver/protocol"
	"google.golang.org/grpc"
)

var log = slog.Scope("RadioClient")

type Callback interface {
	OnData([]complex64)
	OnSmartData([]complex64)
}

type RadioClient struct {
	name           string
	app            string
	address        string
	routineRunning bool
	terminated     bool
	conn           *grpc.ClientConn
	client         protocol.RadioServerClient
	loginData      *protocol.LoginData
	serverInfo     *protocol.ServerInfoData
	deviceInfo     *protocol.DeviceInfo

	currentSampleRate      uint32
	currentSmartSampleRate uint32
	availableSampleRates   []uint32

	iqChannelConfig      *protocol.ChannelConfig
	smartIqChannelConfig *protocol.ChannelConfig

	iqChannelEnabled      bool
	smartIqChannelEnabled bool

	gain      uint32
	streaming bool
	cb        Callback
}

func MakeRadioClient(address, name, application string) *RadioClient {
	return &RadioClient{
		name:                  name,
		app:                   application,
		address:               address,
		routineRunning:        false,
		availableSampleRates:  []uint32{},
		iqChannelConfig:       &protocol.ChannelConfig{},
		smartIqChannelConfig:  &protocol.ChannelConfig{},
		iqChannelEnabled:      false,
		smartIqChannelEnabled: false,
		streaming:             false,
	}
}

// region Public Methods

// GetName returns the name of the active device in RadioClient
func (f *RadioClient) GetName() string {
	if f.deviceInfo != nil {
		return f.deviceInfo.GetDeviceName()
	}

	return "Not Connected"
}

// Start starts the streaming process (if not already started)
func (f *RadioClient) Start() {
	if !f.streaming {
		log.Debug("Starting streaming")
		f.streaming = true
		f.setStreamState()
	}
}

// Stop stop the streaming process (if started)
func (f *RadioClient) Stop() {
	if f.streaming {
		log.Debug("Stopping")
		f.streaming = false
		f.setStreamState()
	}
}

func (f *RadioClient) setStreamState() {
	if f.streaming {
		if f.iqChannelEnabled {
			go f.iqLoop()
		}
		if f.smartIqChannelEnabled {
			go f.smartIqLoop()
		}
	}
}

func (f *RadioClient) smartIqLoop() {
	ctx := context.Background()
	cc := *f.smartIqChannelConfig
	cc.LoginInfo = f.loginData
	iqClient, err := f.client.SmartIQ(ctx, &cc)

	if err != nil {
		log.Fatal(err)
	}
	for f.smartIqChannelEnabled {
		data, err := iqClient.Recv()
		if err != nil {
			log.Error(err)
			f.smartIqChannelEnabled = false
			break
		}
		cData := data.GetComplexSamples()
		if f.cb != nil {
			f.cb.OnSmartData(cData)
		}
	}
}

func (f *RadioClient) iqLoop() {
	ctx := context.Background()
	cc := *f.iqChannelConfig
	cc.LoginInfo = f.loginData
	iqClient, err := f.client.IQ(ctx, &cc)

	if err != nil {
		log.Fatal(err)
	}
	for f.iqChannelEnabled {
		data, err := iqClient.Recv()
		if err != nil {
			log.Error(err)
			f.iqChannelEnabled = false
			break
		}
		cData := data.GetComplexSamples()
		if f.cb != nil {
			f.cb.OnData(cData)
		}
	}
}

// Connect initiates the connection with RadioClient.
// It panics if the connection fails for some reason.
func (f *RadioClient) Connect() {
	if f.routineRunning {
		return
	}

	log.Debug("Trying to connect")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(f.address, opts...)

	if err != nil {
		log.Fatal(err)
	}

	f.conn = conn

	f.client = protocol.NewRadioServerClient(conn)
	ctx := context.Background()

	log.Debug("Connected, sending hello.")
	hr, err := f.client.Hello(ctx, &protocol.HelloData{
		Name:        f.name,
		Application: f.app,
	})

	if err != nil {
		log.Fatal(err)
	}

	f.loginData = hr.Login
	log.Debug("Fetching server info")
	sid, err := f.client.ServerInfo(ctx, &protocol.Empty{})
	if err != nil {
		log.Fatal(err)
	}
	f.serverInfo = sid
	f.deviceInfo = sid.DeviceInfo

	var sampleRates = make([]uint32, f.deviceInfo.MaximumDecimation)
	for i := uint32(0); i < f.deviceInfo.MaximumDecimation; i++ {
		var decim = uint32(1 << i)
		sampleRates[i] = uint32(float32(f.deviceInfo.MaximumSampleRate) / float32(decim))
	}
	f.availableSampleRates = sampleRates
}

// Disconnect disconnects from current connected RadioClient.
func (f *RadioClient) Disconnect() {
	log.Debug("Disconnecting")
	f.terminated = true
	f.iqChannelEnabled = false
	f.smartIqChannelEnabled = false

	if f.conn != nil {
		_ = f.conn.Close()
	}
	f.routineRunning = false
}

// GetSampleRate returns the sample rate of the IQ channel in Hertz
func (f *RadioClient) GetSampleRate() uint32 {
	return f.currentSampleRate
}

// SetSampleRate sets the sample rate of the IQ Channel in Hertz
// Check the available sample rates using GetAvailableSampleRates
// Returns Invalid in case of a invalid value in the input
func (f *RadioClient) SetSampleRate(sampleRate uint32) uint32 {
	for i := uint32(0); i < f.deviceInfo.MaximumDecimation; i++ {
		if f.availableSampleRates[i] == sampleRate {
			f.iqChannelConfig.DecimationStage = i
			f.currentSampleRate = sampleRate
			return sampleRate
		}
	}

	return protocol.Invalid
}

// SetDecimationStage sets the sample rate by using the number of decimation stages.
// Each decimation stage decimates by two, then the total decimation will be defined by 2^stages.
// This is the same as SetSampleRate, but SetSampleRate instead, looks at a pre-filled table of all 2^stages
// decimations that the server supports and applies into the original device sample rate.
func (f *RadioClient) SetDecimationStage(decimation uint32) uint32 {
	if f.deviceInfo == nil || decimation > f.deviceInfo.MaximumDecimation {
		return protocol.Invalid
	}
	f.iqChannelConfig.DecimationStage = decimation
	f.currentSampleRate = f.availableSampleRates[decimation]

	return decimation
}

// GetCenterFrequency returns the IQ Channel Center Frequency in Hz
func (f *RadioClient) GetCenterFrequency() uint32 {
	return f.iqChannelConfig.CenterFrequency
}

// SetCenterFrequency sets the IQ Channel Center Frequency in Hertz and returns it.
func (f *RadioClient) SetCenterFrequency(centerFrequency uint32) uint32 {
	if f.iqChannelConfig.CenterFrequency != centerFrequency {
		f.iqChannelConfig.CenterFrequency = centerFrequency
		if (f.smartIqChannelEnabled) && f.smartIqChannelConfig.CenterFrequency == 0 {
			f.SetSmartCenterFrequency(centerFrequency)
		}
	}

	return f.iqChannelConfig.CenterFrequency
}

// GetSmartCenterFrequency returns the Smart IQ Center Frequency in Hertz
func (f *RadioClient) GetSmartCenterFrequency() uint32 {
	return f.smartIqChannelConfig.CenterFrequency
}

// SetSmartCenterFrequency sets the Smart IQ Center Frequency in Hertz and returns it.
func (f *RadioClient) SetSmartCenterFrequency(centerFrequency uint32) uint32 {
	if f.smartIqChannelConfig.CenterFrequency != centerFrequency {
		f.smartIqChannelConfig.CenterFrequency = centerFrequency
	}

	return f.smartIqChannelConfig.CenterFrequency
}

func (f *RadioClient) SetIQEnabled(iqEnabled bool) {
	f.iqChannelEnabled = iqEnabled
}

func (f *RadioClient) SetSmartIQEnabled(smartIqEnabled bool) {
	f.smartIqChannelEnabled = smartIqEnabled
}

// SetCallback sets the callbacks for server data
func (f *RadioClient) SetCallback(cb Callback) {
	f.cb = cb
}

// GetAvailableSampleRates returns a list of available sample rates for the current connection.
func (f *RadioClient) GetAvailableSampleRates() []uint32 {
	return f.availableSampleRates
}

// SetSmartSampleRate sets the sample rate of the SmartIQ Channel in Hertz
// Check the available sample rates using GetAvailableSampleRates
// Returns Invalid in case of a invalid value in the input
func (f *RadioClient) SetSmartSampleRate(sampleRate uint32) uint32 {
	for i := uint32(0); i < f.deviceInfo.MaximumDecimation; i++ {
		if f.availableSampleRates[i] == sampleRate {
			f.smartIqChannelConfig.DecimationStage = i
			f.currentSmartSampleRate = sampleRate
			return sampleRate
		}
	}

	return protocol.Invalid
}

// SetSmartDecimation sets the sample rate of the Smart IQ by using the number of decimation stages.
// Each decimation stage decimates by two, then the total decimation will be defined by 2^stages.
// This is the same as SetSampleRate, but SetSampleRate instead, looks at a pre-filled table of all 2^stages
// decimations that the server supports and applies into the original device sample rate.
// Returns Invalid in case of a invalid value in the input
func (f *RadioClient) SetSmartDecimation(decimation uint32) uint32 {
	if f.deviceInfo == nil || decimation > f.deviceInfo.MaximumDecimation {
		return protocol.Invalid
	}
	f.smartIqChannelConfig.DecimationStage = decimation
	f.currentSmartSampleRate = f.availableSampleRates[decimation]

	return decimation
}

// GetSmartSampleRate returns the sample rate of Smart IQ in Hertz
func (f *RadioClient) GetSmartSampleRate() uint32 {
	return f.currentSmartSampleRate
}

// SetGain sets the gain stage of the server.
// The actual gain in dB varies from device to device.
// Returns Invalid in case of a invalid value in the input
func (f *RadioClient) SetGain(gain uint32) uint32 {
	if f.deviceInfo == nil || gain > f.deviceInfo.MaximumGain {
		return protocol.Invalid
	}
	f.gain = gain

	return gain
}

// GetGain returns the current gain stage of the server.
func (f *RadioClient) GetGain() uint32 {
	return f.gain
}

// endregion
