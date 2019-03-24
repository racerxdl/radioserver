package tools

import (
	"github.com/racerxdl/segdsp/dsp"
	"math"
)

func Min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func StageToNumber(stage uint32) uint32 {
	return uint32(math.Pow(2, float64(stage)))
}

func GenerateTranslatorTaps(decimation, sampleRate uint32) []float32 {
	var outputSampleRate = float64(sampleRate)
	return dsp.MakeLowPassFixed(1, outputSampleRate, outputSampleRate/(2*float64(decimation)), 31)
}
