package tools

import (
	"bytes"
	"encoding/binary"
	"github.com/racerxdl/segdsp/dsp"
	"math"
)

func Min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func StructToBytes(s interface{}) []uint8 {
	var buff = new(bytes.Buffer)

	_ = binary.Write(buff, binary.LittleEndian, s)

	return buff.Bytes()
}

// region Array Type to Byte Converter
func UnknownArrayToBytes(s interface{}) []uint8 {
	var buff = new(bytes.Buffer)

	_ = binary.Write(buff, binary.LittleEndian, s)

	return buff.Bytes()
}

func Float32ArrayToBytes(s []float32) []uint8 {
	var buff = new(bytes.Buffer)

	for i := 0; i < len(s); i++ {
		_ = binary.Write(buff, binary.LittleEndian, s[i])
	}

	return buff.Bytes()
}

func Float64ArrayToBytes(s []float64) []uint8 {
	var buff = new(bytes.Buffer)

	for i := 0; i < len(s); i++ {
		_ = binary.Write(buff, binary.LittleEndian, s[i])
	}

	return buff.Bytes()
}

func Int16ArrayToBytes(s []int16) []uint8 {
	var buff = new(bytes.Buffer)

	for i := 0; i < len(s); i++ {
		_ = binary.Write(buff, binary.LittleEndian, s[i])
	}

	return buff.Bytes()
}

func Int8ArrayToBytes(s interface{}) []uint8 {
	var buff = new(bytes.Buffer)

	_ = binary.Write(buff, binary.LittleEndian, s)

	return buff.Bytes()
}

func Complex64ArrayToBytes(s []complex64) []uint8 {
	var buff = new(bytes.Buffer)

	for i := 0; i < len(s); i++ {
		_ = binary.Write(buff, binary.LittleEndian, s[i])
	}

	return buff.Bytes()
}

func UInt8ArrayToBytes(s []uint8) []uint8 {
	var buff = make([]uint8, len(s))

	for i := 0; i < len(s); i++ {
		buff[i] = s[i]
	}

	return buff
}
func ArrayToBytes(s interface{}) []uint8 {
	switch v := s.(type) {
	case float32:
		return Float32ArrayToBytes(s.([]float32))
	case float64:
		return Float64ArrayToBytes(s.([]float64))
	case complex64:
		return Complex64ArrayToBytes(s.([]complex64))
	case uint8:
		return UInt8ArrayToBytes(s.([]uint8))
	case int8:
		return Int8ArrayToBytes(s.([]int8))
	case int16:
		return Int16ArrayToBytes(s.([]int16))
	default:
		_ = v
		return UnknownArrayToBytes(s)
	}
}

// endregion
// region Complex64 to XX Array converters
func Complex64ToInt16(samples []complex64) []int16 {
	var i16samples = make([]int16, len(samples)*2)
	for i, v := range samples {
		i16samples[i*2] = int16(real(v) * 32768)
		i16samples[i*2+1] = int16(imag(v) * 32768)
	}
	return i16samples
}

func Complex64ToUInt8(samples []complex64) []uint8 {
	var u8samples = make([]uint8, len(samples)*2)
	for i, v := range samples {
		u8samples[i*2] = uint8(real(v)*127) + 127
		u8samples[i*2+1] = uint8(imag(v)*127) + 127
	}
	return u8samples
}

// endregion
// region Float32 to XX Array converters
func Float32ToInt16(samples []float32) []int16 {
	var i16samples = make([]int16, len(samples))
	for i, v := range samples {
		i16samples[i] = int16(v * 32768)
	}
	return i16samples
}

func Float32ToUInt8(samples []float32) []uint8 {
	var u8samples = make([]uint8, len(samples))
	for i, v := range samples {
		u8samples[i] = uint8(v * 127)
	}
	return u8samples
}

// endregion

func StageToNumber(stage uint32) uint32 {
	return uint32(math.Pow(2, float64(stage)))
}

func GenerateTranslatorTaps(decimation, sampleRate uint32) []float32 {
	var outputSampleRate = float64(sampleRate)
	return dsp.MakeLowPassFixed(1, outputSampleRate, outputSampleRate/(2*float64(decimation)), 31)
}
