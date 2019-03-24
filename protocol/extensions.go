package protocol

import (
	"sync"
	"time"
)

func (m *IQData) GetComplexSamples() []complex64 {
	if m != nil {
		v := make([]complex64, len(m.Samples)/2)
		for i := range v {
			v[i] = complex(m.Samples[i*2], m.Samples[i*2+1])
		}
		return v
	}
	return nil
}

func MakeIQData(channelType ChannelType, samples []complex64) *IQData {
	v := make([]float32, len(samples)*2)

	for i, c := range samples {
		v[i*2] = real(c)
		v[i*2+1] = imag(c)
	}

	return &IQData{
		Timestamp: uint64(time.Now().UnixNano()),
		Status:    StatusType_OK,
		Error:     "",
		Type:      channelType,
		Samples:   v,
	}
}

func MakeIQDataWithPool(channelType ChannelType, samples []complex64, pool sync.Pool) *IQData {
	v := pool.New().([]float32)
	if len(v) < len(samples)*2 {
		v = make([]float32, len(samples)*2)
	} else if len(v) > len(samples)*2 {
		// We dont need to discard, just trim
		v = v[:len(samples)*2]
	}

	for i, c := range samples {
		v[i*2] = real(c)
		v[i*2+1] = imag(c)
	}

	return &IQData{
		Timestamp: uint64(time.Now().UnixNano()),
		Status:    StatusType_OK,
		Error:     "",
		Type:      channelType,
		Samples:   v,
	}
}
