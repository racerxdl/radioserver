package protocol

import (
	"bytes"
	"encoding/binary"
)

func ParseCmdHelloBody(data []uint8) (VersionData, string) {
	var protocolVersion uint64
	var clientName string
	buf := bytes.NewReader(data)
	_ = binary.Read(buf, binary.LittleEndian, &protocolVersion)

	clientName = string(data[8:])

	return SplitProtocolVersion(protocolVersion), clientName
}

func ParseCmdGetSettingBody(data []uint8) {
	// TODO: Implement-me
}

func ParseCmdPingBody(data []uint8) int64 {
	var timestamp int64

	buf := bytes.NewReader(data)
	_ = binary.Read(buf, binary.LittleEndian, &timestamp)

	return timestamp
}

func ParseCmdSetSettingBody(data []uint8) (setting uint32, args []uint32) {
	buf := bytes.NewReader(data)

	var numArgs = (len(data) - 4) / 4

	args = make([]uint32, numArgs)

	_ = binary.Read(buf, binary.LittleEndian, &setting)
	for i := 0; i < numArgs; i++ {
		_ = binary.Read(buf, binary.LittleEndian, &args[i])
	}

	return setting, args
}
