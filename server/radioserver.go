package server

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quan-to/slog"
	"github.com/racerxdl/radioserver"
	"github.com/racerxdl/radioserver/frontends"
	"github.com/racerxdl/radioserver/protocol"
	"github.com/racerxdl/radioserver/tlstools"
	"github.com/racerxdl/radioserver/webapp"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"mime"
	"net"
	"net/http"
	"path"
	"sync"
	"time"
)

var log = slog.Scope("RadioServer")

type RadioServer struct {
	serverInfo *protocol.ServerInfoData

	sessions      map[string]*Session
	sessionLock   sync.Mutex
	grpcServer    *grpc.Server
	tlsGrpcServer *grpc.Server

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

func (rs *RadioServer) Listen(gRPCAddress, httpAddress, tlsGRPCAddress, httpsAddress string) error {
	if rs.grpcServer != nil {
		return fmt.Errorf("server already runing")
	}

	lis, err := net.Listen("tcp", gRPCAddress)
	if err != nil {
		return err
	}

	lisTls, err := net.Listen("tcp", tlsGRPCAddress)
	if err != nil {
		lis.Close()
		return err
	}

	lisHttp, err := net.Listen("tcp", httpAddress)
	if err != nil {
		lis.Close()
		lisTls.Close()
		return err
	}

	lisHttps, err := net.Listen("tcp", httpsAddress)
	if err != nil {
		lis.Close()
		lisTls.Close()
		lisHttp.Close()
		return err
	}

	_, cert := tlstools.GenerateHTTPSKeyPair()

	rs.grpcServer = grpc.NewServer()
	rs.tlsGrpcServer = grpc.NewServer(grpc.Creds(credentials.NewServerTLSFromCert(cert)))

	protocol.RegisterRadioServerServer(rs.grpcServer, rs)
	rs.running = true
	go rs.routines()
	go rs.serve(rs.grpcServer, lis)          // gRPC
	go rs.serve(rs.tlsGrpcServer, lisTls)    // gRPC over TLS
	go rs.serveHttp(lisHttp, lisHttps, cert) // gRPC-WEB
	return nil
}

func (rs *RadioServer) serve(grpc *grpc.Server, conn net.Listener) {
	err := grpc.Serve(conn)
	if err != nil {
		log.Error("RPC Error: %s", err)
	}
	conn.Close()
	rs.Stop()
}

func (rs *RadioServer) serveHttp(conn, conntls net.Listener, cert *tls.Certificate) {
	defer conn.Close()
	defer rs.Stop()

	r := mux.NewRouter()
	p := MakeProxyServer(rs.grpcServer)

	p.RegisterURLs(r)

	r.Handle("/metrics", promhttp.Handler())

	files := webapp.AssetNames()

	for _, f := range files {
		urlPath := path.Join("/", f)
		log.Debug("Registering file %s", urlPath)
		r.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {

			data, err := webapp.Asset(urlPath[1:])
			if err != nil {
				w.WriteHeader(500)
				_, _ = w.Write([]byte("Internal Server Error"))
				return
			}

			ext := path.Ext(urlPath)
			mimeType := mime.TypeByExtension(ext)

			if mimeType == "" {
				mimeType = mime.TypeByExtension(".bin")
			}

			w.Header().Add("content-type", mimeType)
			w.WriteHeader(200)
			_, _ = w.Write(data)
		})
	}

	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		data, err := webapp.Asset("index.html")
		if err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte("Internal Server Error"))
			return
		}

		w.WriteHeader(200)
		_, _ = w.Write(data)
	}

	r.HandleFunc("/", indexHandler)
	r.NotFoundHandler = http.HandlerFunc(indexHandler)

	server := &http.Server{}
	server.Handler = r

	config := &tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0] = *cert

	tlsListener := tls.NewListener(conntls, config)

	serverTls := &http.Server{Handler: r}
	err := http2.ConfigureServer(serverTls, &http2.Server{})
	if err != nil {
		log.Error("Error starting HTTP/2 server: %s", err)
		return
	}

	go func() {
		err = serverTls.Serve(tlsListener)
		if err != nil {
			log.Error("RPC Error: %s", err)
			return
		}
	}()

	err = server.Serve(conn)
	if err != nil {
		log.Error("RPC Error: %s", err)
		return
	}
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
