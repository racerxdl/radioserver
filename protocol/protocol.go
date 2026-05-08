package protocol

import (
	"fmt"
)

func (v *VersionData) AsString() string {
	return fmt.Sprintf("%d.%d - %08x", v.Major, v.Minor, v.Hash)
}

func (v *VersionData) ToUint64() uint64 {
	return GenProtocolVersion(v)
}

func GenProtocolVersion(version *VersionData) uint64 {
	return uint64(((uint64(version.Major)) << 40) | ((uint64(version.Minor)) << 32) | (uint64(version.Hash)))
}

func SplitProtocolVersion(protocol uint64) *VersionData {
	major := uint32(((protocol & (0xFF << 40)) >> 40) & 0xFF)
	minor := uint32(((protocol & (0xFF << 32)) >> 32) & 0xFF)
	hash := uint32(protocol & 0xFFFFFFFF)

	return &VersionData{
		Major: major,
		Minor: minor,
		Hash:  hash,
	}
}

var CurrentProtocolVersion = &VersionData{
	Major: 0,
	Minor: 1,
	Hash:  0,
}

const DefaultPort = 4050

const (
	StatusInvalid = iota
	StatusOK
	StatusError
)

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
