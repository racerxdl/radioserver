package main

import (
	"github.com/racerxdl/radioserver/StateModels"
	"github.com/racerxdl/radioserver/protocol"
	"time"
)

func RunCmdHello(state *StateModels.ClientState) {
	version, name := protocol.ParseCmdHelloBody(state.MessageBody[1:])
	state.Info("Received Hello: %s - %s", version.String(), name)
	state.Name = name
	state.ClientVersion = version

	state.Lock()
	defer state.Unlock()

	data := StateModels.CreateDeviceInfo(state)
	if !state.SendData(data) {
		state.Error("Error sending deviceInfo packet")
	}

	state.SendSync()
}

func RunCmdGetSetting(state *StateModels.ClientState) {
	// TODO
	state.Warn("!!!! RunCmdGetSetting not implemented !!!!")
}

func RunCmdSetSetting(state *StateModels.ClientState) {
	setting, args := protocol.ParseCmdSetSettingBody(state.MessageBody[1:])

	settingName := protocol.SettingNames[setting]

	if !protocol.IsSettingPossible(setting) {
		state.Error("Invalid Setting [%s] => (%d)", settingName, setting)
		return
	}

	state.Debug("Set Setting: %s => %d", settingName, args)

	currentStreaming := state.CGS.Streaming

	if !state.SetSetting(setting, args) {
		return
	}

	state.Lock()
	defer state.Unlock()

	if currentStreaming || currentStreaming != state.CGS.Streaming {
		state.CG.UpdateSettings(state)
		state.SendSync()
	}

	if protocol.SettingAffectsGlobal(setting) {
		serverState.SendSync()
	}
}

func RunCmdPing(state *StateModels.ClientState) {
	timestamp := protocol.ParseCmdPingBody(state.MessageBody[1:])
	delta := float64(time.Now().UnixNano()-timestamp) / 1e6
	state.Debug("Received PING %.2f ms", delta)

	state.Lock()
	defer state.Unlock()

	state.LastPingTime = timestamp
	state.SendPong()
}
