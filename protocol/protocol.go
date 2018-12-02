package protocol

import (
	"fmt"
	"unsafe"
)

type Version struct {
	Major uint8
	Minor uint8
	Hash  uint32
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d - %08x", v.Major, v.Minor, v.Hash)
}

func (v *Version) ToUint64() uint64 {
	return GenProtocolVersion(*v)
}

func GenProtocolVersion(version Version) uint64 {
	return uint64(((uint64(version.Major)) << 40) | ((uint64(version.Minor)) << 32) | (uint64(version.Hash)))
}

func SplitProtocolVersion(protocol uint64) Version {
	major := uint8(((protocol & (0xFF << 40)) >> 40) & 0xFF)
	minor := uint8(((protocol & (0xFF << 32)) >> 32) & 0xFF)
	hash := uint32(protocol & 0xFFFFFFFF)

	return Version{
		Major: major,
		Minor: minor,
		Hash:  hash,
	}
}

const DefaultPort = 4050

// region Internal States
const (
	ParserAcquiringHeader = iota
	ParserReadingData
)

// endregion

// DeviceIds IDs
const (
	DeviceInvalid = iota
	DeviceAirspyOne
	DeviceAirspyHf
	DeviceRtlsdr

	// Radio Server Standard
	DeviceLimeSDRUSB
	DeviceLimeSDRMini
	DeviceSpyServer
	DeviceHackRF
)

// DeviceNames names of the device
const (
	DeviceInvalidName     = "Invalid Device"
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
	DeviceAirspyOne:   DeviceAirspyOneName,
	DeviceAirspyHf:    DeviceAirspyHFName,
	DeviceRtlsdr:      DeviceRtlsdrName,
	DeviceLimeSDRUSB:  DeviceLimeSDRUSBName,
	DeviceLimeSDRMini: DeviceLimeSDRMiniName,
	DeviceHackRF:      DeviceHackRFName,
	DeviceSpyServer:   DeviceSpyserverName,
}

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
	SettingIqFormat
	SettingIqFrequency
	SettingIqDecimation
	SettingDigitalGain
	SettingSmartFrequency
	SettingSmartDecimation
)

// SettingNames list of device names by their ids
var SettingNames = map[uint32]string{
	SettingStreamingMode:    "Streaming Mode",
	SettingStreamingEnabled: "Streaming Enabled",
	SettingGains:            "Gain",
	SettingDigitalGain:      "Digital Gain",
	SettingIqFormat:         "IQ Format",
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

const (
	StreamTypeStatus = iota
	StreamTypeIQ
	StreamTypeSmartIQ
	StreamTypeCombined
)

const (
	StreamFormatInvalid = iota
	StreamFormatInt16
	StreamFormatFloat
)

const (
	MsgTypeDeviceInfo = iota
	MsgTypeClientSync
	MsgTypePong
	MsgTypeReadSetting
	MsgTypeIQ
	MsgTypeSmartIQ
)

type MessageHeader struct {
	PacketNumber    uint32
	ProtocolVersion uint64
	MessageType     uint32
	BodySize        uint32
}

type CommandHeader struct {
	CommandType uint32
	BodySize    uint32
}

type DeviceInfo struct {
	DeviceType        uint32
	DeviceSerial      uint32
	DeviceName        [16]uint8
	MaximumSampleRate uint32
	DecimationStages  uint32
	MaximumGainValue  uint32
	MinimumFrequency  uint32
	MaximumFrequency  uint32
	Resolution        uint32
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

const MessageHeaderSize = uint32(unsafe.Sizeof(MessageHeader{}))
const CommandHeaderSize = uint32(unsafe.Sizeof(CommandHeader{}))
const MaxMessageBodySize = 1 << 20
