//go:build !limesdr && !airspy

package main

import "github.com/racerxdl/radioserver/frontends"

func createFrontend() frontends.Frontend {
	return frontends.CreateTestSignalFrontend()
}
