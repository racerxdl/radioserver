package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/quan-to/slog"
	"github.com/racerxdl/radioserver"
	"github.com/racerxdl/radioserver/protocol"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
	"math"
	"net"
	"time"
)

var log = slog.Scope("RadioClient")
var server = kingpin.Arg("server", "Server Address").Required().String()
var count = kingpin.Flag("c", "Number of pings").Default("5").Int()
var port = kingpin.Flag("port", "Server Port").Default("4050").Int()

var serverIp string

func PingPongTest(client protocol.RadioServerClient) {
	ctx := context.Background()
	sum := int64(0)
	min := int64(math.MaxInt64)
	max := int64(math.MinInt64)
	for i := 0; i < *count; i++ {
		tt := time.Now()
		pd := &protocol.PingData{
			Timestamp: uint64(tt.UnixNano()),
		}
		pong, err := client.Ping(ctx, pd)

		if err != nil {
			log.Fatal(err)
		}

		b, _ := proto.Marshal(pong)

		delta := time.Since(tt).Nanoseconds()
		fmt.Printf("%d bytes from %s (%s): icmp_seq=%d time=%s\n", len(b), *server, serverIp, i+1, time.Duration(delta))
		sum += delta
		if delta < min {
			min = delta
		}
		if delta > max {
			max = delta
		}
		time.Sleep(time.Second)
	}
	avg := sum / int64(*count)
	fmt.Println()
	fmt.Printf("--- %s ping statistics ---\n", *server)
	fmt.Printf("%d packets transmitted, %d received, time %s\n", *count, *count, time.Duration(sum))
	fmt.Printf("rtt min/avg/max/mdev = %s/%s/%s/%s\n", time.Duration(min), time.Duration(avg), time.Duration(max), time.Duration(max-min))
}

func main() {
	kingpin.Version(radioserver.ServerVersion.AsString())

	kingpin.Parse()
	fmt.Printf("--- RadioServer (%s) PING tool ---\n", radioserver.ServerVersion.AsString())

	addr := fmt.Sprintf("%s:%d", *server, *port)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		log.Fatal(err)
	}

	ips, _ := net.LookupIP(*server)

	if len(ips) > 0 {
		serverIp = ips[0].String()
	} else {
		serverIp = *server
	}

	defer conn.Close()
	client := protocol.NewRadioServerClient(conn)

	d, _ := proto.Marshal(&protocol.PingData{
		Timestamp: uint64(time.Now().UnixNano()),
	})
	fmt.Printf("PING %s (%s) %d bytes of data.\n", *server, serverIp, len(d))
	PingPongTest(client)
}
