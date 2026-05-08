//go:build rtlsdr

package main

import "github.com/racerxdl/radioserver/frontends"

func createFrontend() frontends.Frontend {
	return frontends.CreateRTLSDRFrontend(0)
}
