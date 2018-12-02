package main

import (
	"flag"
	"github.com/racerxdl/radioserver/protocol"
)

var ServerVersion = protocol.Version{
	Major: 0,
	Minor: 1,
	Hash:  0,
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
