package server

import (
	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"net/http"
)

func allowAll(_ string) bool {
	return true
}

type ProxyServer struct {
	wrapped *grpcweb.WrappedGrpcServer
	server  *grpc.Server
}

func MakeProxyServer(server *grpc.Server) *ProxyServer {
	return &ProxyServer{
		wrapped: grpcweb.WrapServer(server,
			grpcweb.WithCorsForRegisteredEndpointsOnly(true),
			grpcweb.WithOriginFunc(allowAll)),
		server: server,
	}
}

func (p *ProxyServer) RegisterURLs(r *mux.Router) {
	resources := grpcweb.ListGRPCResources(p.server)
	for _, resource := range resources {
		r.HandleFunc(resource, p.ServeHTTP)
	}
}

func (p *ProxyServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	p.wrapped.ServeHTTP(resp, req)
}

func (p *ProxyServer) IsGrpcWebRequest(req *http.Request) bool {
	return p.wrapped.IsGrpcWebRequest(req)
}

func (p *ProxyServer) IsGrpcWebsocketRequest(req *http.Request) bool {
	return p.wrapped.IsGrpcWebSocketRequest(req)
}
