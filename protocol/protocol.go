package protocol

import (
	"encoding/binary"
	"fmt"
)

func (v *VersionData) AsString() string {
	return fmt.Sprintf("%d.%d - %08x", v.Major, v.Minor, v.Hash)
}

func (v *VersionData) ToUint64() uint64 {
	return GenProtocolVersion(*v)
}

func GenProtocolVersion(version VersionData) uint64 {
	return uint64(((uint64(version.Major)) << 40) | ((uint64(version.Minor)) << 32) | (uint64(version.Hash)))
}

func SplitProtocolVersion(protocol uint64) VersionData {
	major := uint32(((protocol & (0xFF << 40)) >> 40) & 0xFF)
	minor := uint32(((protocol & (0xFF << 32)) >> 32) & 0xFF)
	hash := uint32(protocol & 0xFFFFFFFF)

	return VersionData{
		Major: major,
		Minor: minor,
		Hash:  hash,
	}
}

var CurrentProtocolVersion = VersionData{
	Major: 0,
	Minor: 1,
	Hash:  0,
}

const DefaultPort = 4050

// region Internal States
const (
	GettingHeader = iota
	ReadingData
)

// endregion

const (
	Invalid = iota
	OK
	Error
)

// DeviceIds IDs
const (
	DeviceInvalid = iota
	DeviceTestSignal
	DeviceAirspyOne
	DeviceAirspyHf
	DeviceRtlsdr
	DeviceLimeSDRUSB
	DeviceLimeSDRMini
	DeviceSpyServer
	DeviceHackRF
)

const (
	TypeNone = iota
	TypeDeviceInfo
	TypeClientSync
	TypePong
	TypeReadSetting
	TypeIQ
	TypeSmartIQ
	TypeCombined
	TypeCommand
)

const (
	CmdHello      = 0
	CmdGetSetting = 1
	CmdSetSetting = 2
	CmdPing       = 3
)

const (
	SettingStreamingMode = iota
	SettingStreamingEnabled
	SettingGains
	SettingIqFrequency
	SettingIqDecimation
	SettingDigitalGain
	SettingSmartFrequency
	SettingSmartDecimation
)

// DeviceNames names of the device
const (
	DeviceInvalidName     = "Invalid Device"
	DeviceTestSignalName  = "Test Signal Generator"
	DeviceAirspyOneName   = "Airspy Mini / R2"
	DeviceAirspyHFName    = "Airspy HF / HF+"
	DeviceRtlsdrName      = "RTLSDR"
	DeviceLimeSDRUSBName  = "LimeSDR USB"
	DeviceLimeSDRMiniName = "LimeSDR Mini"
	DeviceHackRFName      = "HackRF"
	DeviceSpyserverName   = "SpyServer"
)

// DeviceName list of device names by their ids
var DeviceName = map[uint32]string{
	DeviceInvalid:     DeviceInvalidName,
	DeviceTestSignal:  DeviceTestSignalName,
	DeviceAirspyOne:   DeviceAirspyOneName,
	DeviceAirspyHf:    DeviceAirspyHFName,
	DeviceRtlsdr:      DeviceRtlsdrName,
	DeviceLimeSDRUSB:  DeviceLimeSDRUSBName,
	DeviceLimeSDRMini: DeviceLimeSDRMiniName,
	DeviceHackRF:      DeviceHackRFName,
	DeviceSpyServer:   DeviceSpyserverName,
}

// SettingNames list of device names by their ids
var SettingNames = map[uint32]string{
	SettingStreamingMode:    "Streaming Mode",
	SettingStreamingEnabled: "Streaming Enabled",
	SettingGains:            "Gain",
	SettingDigitalGain:      "Digital Gain",
	SettingIqFrequency:      "IQ Frequency",
	SettingIqDecimation:     "IQ Decimation",
	SettingSmartFrequency:   "Smart Frequency",
	SettingSmartDecimation:  "Smart Decimation",
}

var PossibleSettings = []uint32{
	SettingStreamingMode,
	SettingStreamingEnabled,
	SettingGains,
	SettingDigitalGain,

	SettingIqFrequency,
	SettingIqDecimation,

	SettingSmartFrequency,
	SettingSmartDecimation,
}

var GlobalAffectedSettings = []uint32{
	SettingGains,
}

func IsSettingPossible(setting uint32) bool {
	for _, v := range PossibleSettings {
		if setting == v {
			return true
		}
	}

	return false
}

func SettingAffectsGlobal(setting uint32) bool {
	for _, v := range GlobalAffectedSettings {
		if setting == v {
			return true
		}
	}

	return false
}

type MessageHeader struct {
	PacketNumber    uint32
	ProtocolVersion uint64
	MessageType     uint32
	Reserved        uint32
	BodySize        uint32
}

type ClientSync struct {
	AllowControl             uint32
	Gains                    [3]uint32
	DeviceCenterFrequency    uint32
	IQCenterFrequency        uint32
	SmartCenterFrequency     uint32
	MinimumIQCenterFrequency uint32
	MaximumIQCenterFrequency uint32
	MinimumSmartFrequency    uint32
	MaximumSmartFrequency    uint32
}

type PingPacket struct {
	Timestamp int64
}

type ReadSettingPacket struct {
	Setting  uint32
	BodySize uint32
}

var MessageHeaderSize = uint32(binary.Size(MessageHeader{}))

const MaxMessageBodySize = 1 << 20
