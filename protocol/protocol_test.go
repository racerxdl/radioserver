package protocol

import "testing"

func TestVersionDataAsString(t *testing.T) {
	v := &VersionData{Major: 1, Minor: 2, Hash: 0xDEADBEEF}
	expected := "1.2 - deadbeef"
	if v.AsString() != expected {
		t.Errorf("expected %q, got %q", expected, v.AsString())
	}
}

func TestGenAndSplitProtocolVersion(t *testing.T) {
	v := &VersionData{Major: 2, Minor: 5, Hash: 0xABCDEF01}
	encoded := GenProtocolVersion(v)
	decoded := SplitProtocolVersion(encoded)

	if decoded.Major != v.Major {
		t.Errorf("Major: expected %d, got %d", v.Major, decoded.Major)
	}
	if decoded.Minor != v.Minor {
		t.Errorf("Minor: expected %d, got %d", v.Minor, decoded.Minor)
	}
	if decoded.Hash != v.Hash {
		t.Errorf("Hash: expected %08x, got %08x", v.Hash, decoded.Hash)
	}
}

func TestVersionDataToUint64(t *testing.T) {
	v := &VersionData{Major: 0, Minor: 1, Hash: 0}
	n := v.ToUint64()
	if n != GenProtocolVersion(v) {
		t.Errorf("ToUint64 mismatch")
	}
}
