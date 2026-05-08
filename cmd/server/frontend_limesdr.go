//go:build limesdr

package main

import "github.com/racerxdl/radioserver/frontends"

func createFrontend() frontends.Frontend {
	return frontends.CreateLimeSDRFrontend(0)
}
