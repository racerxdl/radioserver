package protocol

import (
	"math"
	"sync"
	"testing"
)

func TestGetComplexSamples(t *testing.T) {
	samples := []float32{1.0, 2.0, -3.0, 4.5, 0.0, -1.0}
	iq := &IQData{Samples: samples}

	result := iq.GetComplexSamples()

	if len(result) != 3 {
		t.Fatalf("expected 3 complex samples, got %d", len(result))
	}
	if real(result[0]) != 1.0 || imag(result[0]) != 2.0 {
		t.Errorf("sample 0: expected 1+2i, got %v", result[0])
	}
	if real(result[1]) != -3.0 || imag(result[1]) != 4.5 {
		t.Errorf("sample 1: expected -3+4.5i, got %v", result[1])
	}
	if real(result[2]) != 0.0 || imag(result[2]) != -1.0 {
		t.Errorf("sample 2: expected 0-1i, got %v", result[2])
	}
}

func TestGetComplexSamplesNil(t *testing.T) {
	var iq *IQData
	result := iq.GetComplexSamples()
	if result != nil {
		t.Errorf("expected nil for nil IQData")
	}
}

func TestMakeIQData(t *testing.T) {
	input := []complex64{complex(float32(1.0), float32(2.0)), complex(float32(-3.0), float32(4.0))}
	result := MakeIQData(ChannelType_IQ, input)

	if result.Type != ChannelType_IQ {
		t.Errorf("expected IQ type")
	}
	if result.Status != StatusType_OK {
		t.Errorf("expected OK status")
	}
	if len(result.Samples) != 4 {
		t.Fatalf("expected 4 floats, got %d", len(result.Samples))
	}
	if result.Samples[0] != 1.0 || result.Samples[1] != 2.0 {
		t.Errorf("first sample: expected [1, 2], got [%v, %v]", result.Samples[0], result.Samples[1])
	}
	if result.Samples[2] != -3.0 || result.Samples[3] != 4.0 {
		t.Errorf("second sample: expected [-3, 4], got [%v, %v]", result.Samples[2], result.Samples[3])
	}
}

func TestMakeIQDataRoundtrip(t *testing.T) {
	input := []complex64{
		complex(float32(0.5), float32(-0.5)),
		complex(float32(1.0), float32(0.0)),
		complex(float32(-1.0), float32(-1.0)),
	}
	pb := MakeIQData(ChannelType_SmartIQ, input)
	result := pb.GetComplexSamples()

	if len(result) != len(input) {
		t.Fatalf("length mismatch: %d vs %d", len(result), len(input))
	}
	for i, s := range result {
		if math.Abs(float64(real(s)-real(input[i]))) > 1e-6 {
			t.Errorf("sample %d real: expected %v, got %v", i, real(input[i]), real(s))
		}
		if math.Abs(float64(imag(s)-imag(input[i]))) > 1e-6 {
			t.Errorf("sample %d imag: expected %v, got %v", i, imag(input[i]), imag(s))
		}
	}
}

func TestMakeIQDataWithPool(t *testing.T) {
	input := []complex64{complex(float32(1.0), float32(2.0))}
	pool := &sync.Pool{
		New: func() interface{} {
			s := make([]float32, 4)
			return &s
		},
	}

	result := MakeIQDataWithPool(ChannelType_IQ, input, pool)
	if len(result.Samples) != 2 {
		t.Fatalf("expected 2 floats, got %d", len(result.Samples))
	}
	if result.Samples[0] != 1.0 || result.Samples[1] != 2.0 {
		t.Errorf("expected [1, 2], got %v", result.Samples)
	}
}
