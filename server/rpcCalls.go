package server

import (
	"context"
	"fmt"
	"github.com/racerxdl/radioserver/protocol"
	"runtime"
	"sync"
	"time"
)

// region GRPC Stuff
func (rs *RadioServer) Me(ctx context.Context, data *protocol.LoginData) (*protocol.MeData, error) {
	token := data.Token
	session, ok := rs.sessions[token]
	if ok {
		return &session.MeData, nil
	}
	return nil, fmt.Errorf("not logged in")
}

func (rs *RadioServer) Hello(ctx context.Context, hdata *protocol.HelloData) (*protocol.HelloReturn, error) {
	rs.sessionLock.Lock()
	defer rs.sessionLock.Unlock()

	s := GenerateSession(hdata, false)
	rs.sessions[s.ID] = s
	log.Info("Welcome %s!", s.Name)

	return &protocol.HelloReturn{
		Status: protocol.OK,
		Login:  &s.LoginData,
	}, nil
}

func (rs *RadioServer) Bye(ctx context.Context, ld *protocol.LoginData) (*protocol.ByeReturn, error) {
	rs.sessionLock.Lock()
	defer rs.sessionLock.Unlock()

	s := rs.sessions[ld.Token]
	if s == nil {
		return nil, fmt.Errorf("not logged in")
	}

	delete(rs.sessions, ld.Token)

	s.FullStop()

	log.Info("Bye %s!", s.Name)
	return &protocol.ByeReturn{
		Message: "Farewell, my friend!",
	}, nil
}

func (rs *RadioServer) ServerInfo(context.Context, *protocol.Empty) (*protocol.ServerInfoData, error) {
	return rs.serverInfo, nil
}

func (rs *RadioServer) SmartIQ(cc *protocol.ChannelConfig, server protocol.RadioServer_SmartIQServer) error {
	rs.sessionLock.Lock()
	s := rs.sessions[cc.LoginInfo.Token]
	if s == nil {
		return fmt.Errorf("not logged in")
	}
	rs.sessionLock.Unlock()

	if s.CG.SmartIQRunning() {
		return fmt.Errorf("already running")
	}

	s.CG.UpdateSettings(protocol.ChannelType_SmartIQ, rs.frontend, cc)
	s.CG.StartSmartIQ()
	defer s.CG.StopSmartIQ()

	lastNumSamples := 0
	pool := sync.Pool{
		New: func() interface{} {
			return make([]float32, lastNumSamples)
		},
	}

	for {
		for s.SmartIQFifo.Len() > 0 {
			samples := s.SmartIQFifo.Next().([]complex64)
			pb := protocol.MakeIQDataWithPool(protocol.ChannelType_IQ, samples, pool)
			if err := server.Send(pb); err != nil {
				log.Error("Error sending samples to %s: %s", s.Name, err)
				return err
			}
			s.KeepAlive()

			if len(pb.Samples) != lastNumSamples {
				lastNumSamples = len(pb.Samples)
			}

			pool.Put(pb.Samples) // If the size is not correct, MakeIQDataWithPool will discard or trim it

			if s.IsFullStopped() {
				log.Error("Session Expired")
				return fmt.Errorf("session expired")
			}
			runtime.Gosched()
		}
		time.Sleep(time.Millisecond)
	}
}

func (rs *RadioServer) IQ(cc *protocol.ChannelConfig, server protocol.RadioServer_IQServer) error {
	rs.sessionLock.Lock()
	s := rs.sessions[cc.LoginInfo.Token]
	if s == nil {
		return fmt.Errorf("not logged in")
	}
	rs.sessionLock.Unlock()

	if s.CG.IQRunning() {
		return fmt.Errorf("already running")
	}

	s.CG.UpdateSettings(protocol.ChannelType_IQ, rs.frontend, cc)
	s.CG.StartIQ()
	defer s.CG.StopIQ()

	lastNumSamples := 0
	pool := sync.Pool{
		New: func() interface{} {
			return make([]float32, lastNumSamples)
		},
	}

	for {
		for s.IQFifo.Len() > 0 {
			samples := s.IQFifo.Next().([]complex64)
			pb := protocol.MakeIQDataWithPool(protocol.ChannelType_IQ, samples, pool)
			if err := server.Send(pb); err != nil {
				log.Error("Error sending samples to %s: %s", s.Name, err)
				return err
			}
			s.KeepAlive()

			if len(pb.Samples) != lastNumSamples {
				lastNumSamples = len(pb.Samples)
			}

			pool.Put(pb.Samples) // If the size is not correct, MakeIQDataWithPool will discard or trim it

			if s.IsFullStopped() {
				log.Error("Session Expired")
				return fmt.Errorf("session expired")
			}
			runtime.Gosched()
		}
		time.Sleep(time.Millisecond)
	}
}

func (rs *RadioServer) Ping(ctx context.Context, pd *protocol.PingData) (*protocol.PingData, error) {
	if pd.Token != "" {
		rs.sessionLock.Lock()
		session := rs.sessions[pd.Token]
		if session != nil {
			session.KeepAlive()
		}
		rs.sessionLock.Unlock()
	}

	return &protocol.PingData{
		Timestamp: uint64(time.Now().UnixNano()),
	}, nil
}

// endregion
