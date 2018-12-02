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
		if state.CurrentState == protocol.ParserAcquiringHeader {
			for state.CurrentState == protocol.ParserAcquiringHeader && len(buffer) > 0 {
				consumed = parseHeader(state, buffer)
				buffer = buffer[consumed:]
			}

			if state.CurrentState == protocol.ParserReadingData {

				if state.Cmd.BodySize > protocol.MaxMessageBodySize {
					state.Error("Client sent an BodySize of %d which is higher than max %d", state.Cmd.BodySize, protocol.MaxMessageBodySize)
					state.Running = false
					return
				}

				state.CmdBody = make([]uint8, state.Cmd.BodySize)
			}
		}

		if state.CurrentState == protocol.ParserReadingData {
			consumed = parseBody(state, buffer)
			buffer = buffer[consumed:]

			if state.CurrentState == protocol.ParserAcquiringHeader {
				state.CmdReceived++
				runCommand(state)
			}
		}
	}
}

func parseBody(state *StateModels.ClientState, buffer []uint8) uint32 {
	consumed := uint32(0)

	for len(buffer) > 0 {
		toWrite := tools.Min(state.Cmd.BodySize-state.ParserPosition, uint32(len(buffer)))
		for i := uint32(0); i < toWrite; i++ {
			state.CmdBody[i+state.ParserPosition] = buffer[i]
		}
		buffer = buffer[toWrite:]
		consumed += toWrite
		state.ParserPosition += toWrite

		if state.ParserPosition == state.Cmd.BodySize {
			state.ParserPosition = 0
			state.CurrentState = protocol.ParserAcquiringHeader
			return consumed
		}
	}

	return consumed
}

func parseHeader(state *StateModels.ClientState, buffer []uint8) uint32 {
	consumed := uint32(0)

	for len(buffer) > 0 {
		toWrite := tools.Min(protocol.CommandHeaderSize-state.ParserPosition, uint32(len(buffer)))
		for i := uint32(0); i < toWrite; i++ {
			state.HeaderBuffer[i+state.ParserPosition] = buffer[i]
		}
		buffer = buffer[toWrite:]
		consumed += toWrite
		state.ParserPosition += toWrite

		if state.ParserPosition == protocol.CommandHeaderSize {
			state.ParserPosition = 0
			buf := bytes.NewReader(state.HeaderBuffer)
			err := binary.Read(buf, binary.LittleEndian, &state.Cmd)
			if err != nil {
				panic(err)
			}

			if state.Cmd.BodySize > 0 {
				state.CurrentState = protocol.ParserReadingData
			}

			return consumed
		}
	}

	return consumed
}

func runCommand(state *StateModels.ClientState) {
	var cmdType = state.Cmd.CommandType

	if cmdType == protocol.CmdHello {
		RunCmdHello(state)
	} else if cmdType == protocol.CmdGetSetting {
		RunCmdGetSetting(state)
	} else if cmdType == protocol.CmdSetSetting {
		RunCmdSetSetting(state)
	} else if cmdType == protocol.CmdPing {
		RunCmdPing(state)
	}
}
