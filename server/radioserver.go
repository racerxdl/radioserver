package server

import (
	"fmt"
	"github.com/quan-to/slog"
	"github.com/racerxdl/radioserver"
	"github.com/racerxdl/radioserver/frontends"
	"github.com/racerxdl/radioserver/protocol"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

var log = slog.Scope("RadioServer")

type RadioServer struct {
	serverInfo *protocol.ServerInfoData

	sessions    map[string]*Session
	sessionLock sync.Mutex
	grpcServer  *grpc.Server

	frontend frontends.Frontend
	running  bool

	lastSessionChecks time.Time
}

func MakeRadioServer(frontend frontends.Frontend) *RadioServer {
	rs := &RadioServer{
		serverInfo: &protocol.ServerInfoData{
			ControlAllowed:           false,
			ServerCenterFrequency:    0,
			MinimumIQCenterFrequency: 0,
			MaximumIQCenterFrequency: 0,
			MinimumSmartFrequency:    0,
			MaximumSmartFrequency:    0,
			DeviceInfo: &protocol.DeviceInfo{
				DeviceType:        frontend.GetDeviceType(),
				DeviceSerial:      frontend.GetDeviceSerial(),
				DeviceName:        frontend.GetName(),
				MaximumSampleRate: frontend.GetMaximumSampleRate(),
				MaximumGain:       frontend.MaximumGainValue(),
				MaximumDecimation: frontend.MaximumDecimationStages(),
				MinimumFrequency:  frontend.MinimumFrequency(),
				MaximumFrequency:  frontend.MaximumFrequency(),
				ADCResolution:     uint32(frontend.GetResolution()),
			},
			Version: &protocol.VersionData{
				Major: uint32(radioserver.ServerVersion.Major),
				Minor: uint32(radioserver.ServerVersion.Minor),
				Hash:  radioserver.ServerVersion.Hash,
			},
		},
		sessions:    map[string]*Session{},
		sessionLock: sync.Mutex{},
		frontend:    frontend,
	}

	frontend.SetSamplesAvailableCallback(rs.onSamples)

	return rs
}

func (rs *RadioServer) Listen(address string) error {
	if rs.grpcServer != nil {
		return fmt.Errorf("server already runing")
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	rs.grpcServer = grpc.NewServer()

	protocol.RegisterRadioServerServer(rs.grpcServer, rs)
	rs.running = true
	go rs.routines()
	go rs.serve(lis)
	return nil
}

func (rs *RadioServer) serve(conn net.Listener) {
	err := rs.grpcServer.Serve(conn)
	if err != nil {
		log.Error("RPC Error: %s", err)
	}
	rs.Stop()
}

func (rs *RadioServer) Stop() {
	if rs.grpcServer == nil {
		return
	}
	log.Info("Stopping RPC Server")
	rs.grpcServer.Stop()
	rs.grpcServer = nil
	rs.running = false
}
