package radioserver

import (
	"github.com/racerxdl/radioserver/protocol"
	"strconv"
)

var ServerVersion = protocol.VersionData{
	Major: 0,
	Minor: 1,
	Hash:  0,
}

func init() {
	hash, _ := strconv.ParseInt(commitHash, 16, 32)
	ServerVersion.Hash = uint32(hash)
}
