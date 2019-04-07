package server

import (
	uuid2 "github.com/gofrs/uuid"
	"github.com/racerxdl/go.fifo"
	"github.com/racerxdl/radioserver/DSP"
	"github.com/racerxdl/radioserver/protocol"
	"time"
)

const (
	expirationTime = time.Second * 120
	maxFifoBuffs   = 4096
)

type Session struct {
	protocol.MeData
	protocol.LoginData

	ID         string
	LastUpdate time.Time

	IQFifo        *fifo.Queue
	SmartIQFifo   *fifo.Queue
	FrequencyFifo *fifo.Queue
	CG            *DSP.ChannelGenerator

	fullStopped bool
}

func GenerateSession(loginData *protocol.HelloData, controlAllowed bool) *Session {
	u, _ := uuid2.NewV4()
	ID := u.String()

	CG := DSP.CreateChannelGenerator()

	s := &Session{
		IQFifo:        fifo.NewQueue(),
		SmartIQFifo:   fifo.NewQueue(),
		FrequencyFifo: fifo.NewQueue(),
		MeData: protocol.MeData{
			Name:           loginData.Name,
			Application:    loginData.Application,
			Username:       loginData.Username,
			ControlAllowed: controlAllowed,
		},
		LoginData: protocol.LoginData{
			Token: ID,
		},
		ID:          ID,
		LastUpdate:  time.Now(),
		CG:          CG,
		fullStopped: false,
	}

	CG.SetOnIQ(func(samples []complex64) {
		if s.IQFifo.Len() < maxFifoBuffs && !s.fullStopped {
			s.IQFifo.Add(samples)
		}
	})

	CG.SetOnSmartIQ(func(samples []complex64) {
		if s.SmartIQFifo.Len() < maxFifoBuffs && !s.fullStopped {
			s.SmartIQFifo.Add(samples)
		}
	})

	CG.SetOnFC(func(samples []float32) {
		if s.FrequencyFifo.Len() < maxFifoBuffs && !s.fullStopped {
			s.FrequencyFifo.Add(samples)
		}
	})

	CG.Start()

	return s
}

func (s *Session) Expired() bool {
	return time.Since(s.LastUpdate) > expirationTime
}

func (s *Session) KeepAlive() {
	s.LastUpdate = time.Now()
}

func (s *Session) IsFullStopped() bool {
	return s.fullStopped
}

func (s *Session) FullStop() {
	s.CG.StopIQ()
	s.CG.StopSmartIQ()
	s.CG.Stop()
	s.fullStopped = true
}
