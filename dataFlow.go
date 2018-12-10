package main

import (
	"bytes"
	"encoding/binary"
	"github.com/racerxdl/radioserver/StateModels"
	"github.com/racerxdl/radioserver/protocol"
	"github.com/racerxdl/radioserver/tools"
)

func parseMessage(state *StateModels.ClientState, buffer []uint8) {
	state.ReceivedBytes += uint64(len(buffer))

	var consumed uint32

	for len(buffer) > 0 && tcpServerStatus && state.Running {
		if state.CurrentState == protocol.GettingHeader {
			for state.CurrentState == protocol.GettingHeader && len(buffer) > 0 {
				consumed = parseHeader(state, buffer)
				buffer = buffer[consumed:]
			}

			if state.CurrentState == protocol.ReadingData {
				if state.Message.BodySize > protocol.MaxMessageBodySize {
					state.Error("Client sent an BodySize of %d which is higher than max %d", state.Message.BodySize, protocol.MaxMessageBodySize)
					state.Running = false
					return
				}
				state.MessageBody = make([]uint8, state.Message.BodySize)
			}
		}

		if state.CurrentState == protocol.ReadingData {
			consumed = parseBody(state, buffer)
			buffer = buffer[consumed:]

			if state.CurrentState == protocol.GettingHeader {
				state.CmdReceived++
				runCommand(state)
			}
		}
	}
}

func parseBody(state *StateModels.ClientState, buffer []uint8) uint32 {
	consumed := uint32(0)

	for len(buffer) > 0 {
		toWrite := tools.Min(state.Message.BodySize-state.ParserPosition, uint32(len(buffer)))
		for i := uint32(0); i < toWrite; i++ {
			state.MessageBody[i+state.ParserPosition] = buffer[i]
		}
		buffer = buffer[toWrite:]
		consumed += toWrite
		state.ParserPosition += toWrite

		if state.ParserPosition == state.Message.BodySize {
			state.ParserPosition = 0
			state.CurrentState = protocol.GettingHeader
			return consumed
		}
	}

	return consumed
}

func parseHeader(state *StateModels.ClientState, buffer []uint8) uint32 {
	consumed := uint32(0)

	for len(buffer) > 0 {
		toWrite := tools.Min(protocol.MessageHeaderSize-state.ParserPosition, uint32(len(buffer)))
		for i := uint32(0); i < toWrite; i++ {
			state.HeaderBuffer[i+state.ParserPosition] = buffer[i]
		}
		buffer = buffer[toWrite:]
		consumed += toWrite
		state.ParserPosition += toWrite

		if state.ParserPosition == protocol.MessageHeaderSize {
			state.ParserPosition = 0
			buf := bytes.NewReader(state.HeaderBuffer)
			err := binary.Read(buf, binary.LittleEndian, &state.Message)
			if err != nil {
				panic(err)
			}

			if state.Message.BodySize > 0 {
				state.CurrentState = protocol.ReadingData
			}

			return consumed
		}
	}

	return consumed
}

func runCommand(state *StateModels.ClientState) {
	var cmdType = state.MessageBody[0]

	if cmdType == protocol.CmdHello {
		RunCmdHello(state)
	} else if cmdType == protocol.CmdGetSetting {
		RunCmdGetSetting(state)
	} else if cmdType == protocol.CmdSetSetting {
		RunCmdSetSetting(state)
	} else if cmdType == protocol.CmdPing {
		RunCmdPing(state)
	} else {
		state.Error("Unknown Command %d", cmdType)
	}
}
