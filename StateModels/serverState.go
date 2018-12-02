package StateModels

import (
	"github.com/racerxdl/radioserver/SLog"
	"github.com/racerxdl/radioserver/frontends"
	"github.com/racerxdl/radioserver/protocol"
	"sync"
)

type ServerState struct {
	DeviceInfo    protocol.DeviceInfo
	clients       []*ClientState
	clientListMtx sync.Mutex
	Frontend      frontends.Frontend
	CanControl    uint32
}

func CreateServerState() *ServerState {
	return &ServerState{
		clientListMtx: sync.Mutex{},
		clients:       make([]*ClientState, 0),
	}
}

func (s *ServerState) indexOfClient(state *ClientState) int {
	for k, v := range s.clients {
		if v.UUID == state.UUID {
			return k
		}
	}

	return -1
}

func (s *ServerState) PushClient(state *ClientState) {
	s.clientListMtx.Lock()
	defer s.clientListMtx.Unlock()

	count := len(s.clients)

	s.clients = append(s.clients, state)
	if count == 0 {
		SLog.Info("First client connected. Starting frontend...")
		s.Frontend.Start()
	}
}

func (s *ServerState) RemoveClient(state *ClientState) {
	s.clientListMtx.Lock()
	defer s.clientListMtx.Unlock()
	idx := s.indexOfClient(state)
	if idx != -1 {
		s.clients = append(s.clients[:idx], s.clients[idx+1:]...)
	}

	if len(s.clients) == 0 {
		SLog.Info("Last client gone. Stopping frontend...")
		s.Frontend.Stop()
	}
}

func (s *ServerState) SendSync() bool {
	s.clientListMtx.Lock()
	defer s.clientListMtx.Unlock()

	for i := 0; i < len(s.clients); i++ {
		s.clients[i].SendSync()
	}

	return true
}

func (s *ServerState) PushSamples(samples []complex64) {
	var clientList []*ClientState
	s.clientListMtx.Lock()
	clientList = make([]*ClientState, len(s.clients))
	copy(clientList, s.clients)
	s.clientListMtx.Unlock()

	for _, v := range clientList {
		v.CG.PushSamples(samples)
	}
}
