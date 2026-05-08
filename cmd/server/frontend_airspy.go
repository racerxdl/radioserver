//go:build airspy

package main

import "github.com/racerxdl/radioserver/frontends"

func createFrontend() frontends.Frontend {
	return frontends.CreateAirspyFrontend(0)
}
