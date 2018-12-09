package StateModels

import (
	"github.com/google/uuid"
	"github.com/racerxdl/radioserver/SLog"
	"github.com/racerxdl/radioserver/protocol"
	"github.com/racerxdl/radioserver/tools"
	"net"
	"strings"
	"sync"
	"time"
)

type ChannelGeneratorState struct {
	Streaming     bool
	StreamingMode uint32

	// Channel Mode
	IQCenterFrequency uint32
	IQDecimation      uint32

	// Smart Settings
	SmartIQDecimation    uint32
	SmartCenterFrequency uint32
}

// region ClientState

type ClientState struct {
	sync.Mutex
	UUID           string
	Buffer         []uint8
	HeaderBuffer   []uint8
	LogInstance    *SLog.Instance
	Addr           net.Addr
	Conn           net.Conn
	Running        bool
	Name           string
	ClientVersion  protocol.Version
	CurrentState   int
	ReceivedBytes  uint64
	SentBytes      uint64
	ConnectedSince time.Time
	CmdReceived    uint64
	SentPackets    uint64

	ServerVersion protocol.Version
	ServerState   *ServerState

	// Command State
	Cmd            protocol.CommandHeader
	CmdBody        []uint8
	ParserPosition uint32
	SyncInfo       protocol.ClientSync

	LastPingTime int64

	// Channel Generator
	CGS ChannelGeneratorState
	CG  *ChannelGenerator
}

func CreateClientState(centerFrequency uint32) *ClientState {
	var cs = &ClientState{
		UUID:           uuid.New().String(),
		Buffer:         make([]uint8, 64*1024),
		CurrentState:   protocol.ParserAcquiringHeader,
		ConnectedSince: time.Now(),
		ReceivedBytes:  0,
		SentBytes:      0,
		Running:        false,
		SentPackets:    0,
		CmdReceived:    0,
		ParserPosition: 0,
		LogInstance:    SLog.Scope("ClientState"),
		HeaderBuffer:   make([]uint8, protocol.MessageHeaderSize),
		CGS: ChannelGeneratorState{
			Streaming:            false,
			StreamingMode:        protocol.StreamTypeIQ,
			IQCenterFrequency:    centerFrequency,
			IQDecimation:         0,
			SmartIQDecimation:    0,
			SmartCenterFrequency: centerFrequency,
		},
		CG: CreateChannelGenerator(),
	}

	cs.CG.SetOnSmartIQ(cs.onSmart)
	cs.CG.SetOnIQ(cs.onIQ)

	return cs
}

func (state *ClientState) Log(str interface{}, v ...interface{}) *ClientState {
	state.LogInstance.Log(str, v...)
	return state
}

func (state *ClientState) Info(str interface{}, v ...interface{}) *ClientState {
	state.LogInstance.Info(str, v...)
	return state
}

func (state *ClientState) Debug(str interface{}, v ...interface{}) *ClientState {
	state.LogInstance.Debug(str, v...)
	return state
}

func (state *ClientState) Warn(str interface{}, v ...interface{}) *ClientState {
	state.LogInstance.Warn(str, v...)
	return state
}

func (state *ClientState) Error(str interface{}, v ...interface{}) *ClientState {
	state.LogInstance.Error(str, v...)
	return state
}

func (state *ClientState) Fatal(str interface{}, v ...interface{}) {
	state.LogInstance.Fatal(str, v)
}

func (state *ClientState) FullStop() {
	state.Info("Fully stopping Client")
	state.CG.Stop()
	state.Info("Client stopped")
}

func (state *ClientState) SendData(buffer []uint8) bool {
	n, err := state.Conn.Write(buffer)
	if err != nil {
		errMsg := err.Error()
		if !strings.Contains(errMsg, "closed") && !strings.Contains(errMsg, "broken pipe") {
			state.LogInstance.Error("Error sending data: %s", err)
		}
		return false
	}

	state.SentPackets++

	if n > 0 {
		state.SentBytes += uint64(n)
	}

	return true
}

func (state *ClientState) onSmart(samples []complex64) {
	samplesToSend := tools.Complex64ToInt16(samples)
	state.Lock()
	defer state.Unlock()

	if samplesToSend != nil {
		var data = CreateDataPacket(state, protocol.MsgTypeSmartIQ, samplesToSend)
		state.SendData(data)
	}
}

func (state *ClientState) onIQ(samples []complex64) {
	samplesToSend := tools.Complex64ToInt16(samples)
	msgType := protocol.MsgTypeIQ

	if samplesToSend != nil {
		state.SendIQ(samplesToSend, uint32(msgType))
	}
}

func (state *ClientState) SendIQ(samples interface{}, messageType uint32) {
	state.Lock()
	defer state.Unlock()

	var bodyData = tools.ArrayToBytes(samples)

	var header = protocol.MessageHeader{
		ProtocolVersion: state.ServerVersion.ToUint64(),
		MessageType:     messageType,
		PacketNumber:    uint32(state.SentPackets & 0xFFFFFFFF),
		BodySize:        uint32(len(bodyData)),
	}

	if len(bodyData) > protocol.MaxMessageBodySize {
		// Segmentation
		for len(bodyData) > 0 {
			chunkSize := tools.Min(protocol.MaxMessageBodySize, uint32(len(bodyData)))
			segment := bodyData[:chunkSize]
			bodyData = bodyData[chunkSize:]
			header.BodySize = uint32(len(segment))
			header.PacketNumber = uint32(state.SentPackets & 0xFFFFFFFF)
			state.SendData(CreateRawPacket(header, segment))
		}
		return
	}

	state.SendData(CreateRawPacket(header, bodyData))
}

func (state *ClientState) updateSync() {
	var halfSampleRate = state.ServerState.Frontend.GetSampleRate() / 2
	var centerFreq = state.CGS.IQCenterFrequency

	state.SyncInfo.SmartCenterFrequency = state.CGS.SmartCenterFrequency
	state.SyncInfo.IQCenterFrequency = state.CGS.IQCenterFrequency
	state.SyncInfo.AllowControl = state.ServerState.CanControl
	state.SyncInfo.Gains = [3]uint32{uint32(state.ServerState.Frontend.GetGain()), 0, 0}
	state.SyncInfo.DeviceCenterFrequency = state.ServerState.Frontend.GetCenterFrequency()
	state.SyncInfo.MaximumIQCenterFrequency = centerFreq + halfSampleRate
	state.SyncInfo.MinimumIQCenterFrequency = centerFreq - halfSampleRate
	state.SyncInfo.MaximumSmartFrequency = centerFreq + halfSampleRate
	state.SyncInfo.MinimumSmartFrequency = centerFreq - halfSampleRate
}

func (state *ClientState) SendSync() {
	state.updateSync()
	data := CreateClientSync(state)
	if !state.SendData(data) {
		state.Error("Error sending syncInfo packet")
	}
}

func (state *ClientState) SendPong() {
	data := CreatePong(state)
	if !state.SendData(data) {
		state.Error("Error sending pong packet")
	}
}

func (state *ClientState) SetSetting(setting uint32, args []uint32) bool {
	switch setting {
	case protocol.SettingStreamingMode:
		return state.SetStreamingMode(args[0])
	case protocol.SettingStreamingEnabled:
		return state.SetStreamingEnabled(args[0] == 1)
	case protocol.SettingGains:
		return state.SetGain(args[0])
	case protocol.SettingIqFrequency:
		return state.SetIQFrequency(args[0])
	case protocol.SettingIqDecimation:
		return state.SetIQDecimation(args[0])
	case protocol.SettingSmartFrequency:
		return state.SetSmartFrequency(args[0])
	case protocol.SettingSmartDecimation:
		return state.SetSmartDecimation(args[0])
	}

	return false
}

func (state *ClientState) SetStreamingMode(mode uint32) bool {
	state.CGS.StreamingMode = mode
	return true
}
func (state *ClientState) SetStreamingEnabled(enabled bool) bool {
	var enabledString = "Enabled"
	if !enabled {
		enabledString = "Disabled"
	}

	state.Log("Streaming %s", enabledString)
	state.CGS.Streaming = enabled

	return true
}
func (state *ClientState) SetGain(gain uint32) bool {
	state.ServerState.Frontend.SetGain(uint8(gain))
	return true
}
func (state *ClientState) SetIQFrequency(frequency uint32) bool {
	state.CGS.IQCenterFrequency = frequency
	state.updateSync()
	return true
}
func (state *ClientState) SetIQDecimation(decimation uint32) bool {
	if state.ServerState.DeviceInfo.DecimationStages >= decimation {
		state.CGS.IQDecimation = decimation
		return true
	}

	return false
}
func (state *ClientState) SetSmartFrequency(frequency uint32) bool {
	state.CGS.SmartCenterFrequency = frequency
	state.updateSync()
	return true
}

func (state *ClientState) SetSmartDecimation(decimation uint32) bool {
	if state.ServerState.DeviceInfo.DecimationStages >= decimation {
		state.CGS.SmartIQDecimation = decimation
		return true
	}

	return false
}

// endregion
