package client

// ComplexInt16 is a Complex Number in a signed 16 bit number
type ComplexInt16 struct {
	Real int16
	Imag int16
}

// ComplexUInt16 is a Complex Number in a unsigned 16 bit number
type ComplexUInt16 struct {
	Real uint16
	Imag uint16
}

// ComplexUInt8 is a Complex Number in a unsigned 8 bit number
// In this case the value 0 is in variable half-way (127)
type ComplexUInt8 struct {
	Real uint8
	Imag uint8
}

// Return types for callback OnData
const (
	SamplesComplex64 = iota
	SmartSamplesComplex64
	SamplesComplex32
	SmartSamplesComplex32
	SamplesBytes
	DeviceSync
)

type Callback interface {
	OnData(int, interface{})
}
