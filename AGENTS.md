# RadioServer

Go gRPC SDR server. Streams IQ samples from SDR hardware to multiple clients with independent per-session DSP.

## Commands

```bash
# Default (TestSignal, no hardware needed)
go build -o radioserver ./cmd/server/

# With LimeSDR (requires LimeSuite headers)
go build -tags limesdr -o radioserver ./cmd/server/

# With Airspy (requires libairspy)
go build -tags airspy -o radioserver ./cmd/server/

go test -v -race ./...
golangci-lint run
go vet ./...
go fmt ./...

# Regenerate protobuf (after changing server.proto)
protoc -I protocol/ protocol/server.proto --go_out=protocol/ --go_opt=paths=source_relative --go-grpc_out=protocol/ --go-grpc_opt=paths=source_relative
```

No Makefile or task runner exists.

## Architecture

- `cmd/server/` — Server binary. Uses build tags (`limesdr`, `airspy`) to select frontend. Default: TestSignal.
- `cmd/client/` — Test client binary
- `cmd/ping/` — gRPC ping/latency tool
- `protocol/` — Protobuf definition + generated code + legacy binary protocol types
- `server/` — gRPC server: session management, RPC handlers, streaming
- `client/` — Client library (used by segdsp)
- `DSP/` — Per-session ChannelGenerator (frequency translation, decimation, Blackman-Harris windowing)
- `frontends/` — SDR hardware abstractions (build-tag isolated):
  - `TestSignal.go` — Software signal generator (default, no hardware)
  - `LimeSDR.go` — `//go:build limesdr` (requires LimeSuite CGO)
  - `Airspy.go` — `//go:build airspy` (requires libairspy CGO)
- `tools/` — Utility functions (filter tap generation)

## Gotchas

- LimeSDR and Airspy frontends require CGO and system libraries; use build tags
- Default build (`go build ./cmd/server/`) uses TestSignal — no hardware needed
- Each connected client gets its own `ChannelGenerator` with independent tuning
- Session expiration: 120 seconds of inactivity
- SmartIQ rate-limited to 20 fps, frame size 4096 samples
- Server listens on `:4050` by default

## Testing

No tests currently exist. All code is integration-level (requires hardware or gRPC server running).

## Conventions

- gRPC service with protobuf, generated code in `protocol/`
- Frontends implement the `frontends.Frontend` interface
- Client library uses `Callback` interface (`OnData`, `OnSmartData`)
- Session management via UUID tokens in `sync.Mutex`-protected map
- FIFO-based sample buffering with `go.fifo`
