package tools

import "testing"

func TestMin(t *testing.T) {
	if Min(3, 5) != 3 {
		t.Error("Min(3,5) should be 3")
	}
	if Min(5, 3) != 3 {
		t.Error("Min(5,3) should be 3")
	}
	if Min(3, 3) != 3 {
		t.Error("Min(3,3) should be 3")
	}
}

func TestStageToNumber(t *testing.T) {
	tests := []struct {
		stage    uint32
		expected uint32
	}{
		{0, 1},
		{1, 2},
		{2, 4},
		{3, 8},
		{4, 16},
	}
	for _, tt := range tests {
		result := StageToNumber(tt.stage)
		if result != tt.expected {
			t.Errorf("StageToNumber(%d) = %d, expected %d", tt.stage, result, tt.expected)
		}
	}
}

func TestGenerateTranslatorTaps(t *testing.T) {
	taps := GenerateTranslatorTaps(2, 1000000)
	if len(taps) != 31 {
		t.Fatalf("expected 31 taps, got %d", len(taps))
	}
	for _, tap := range taps {
		if tap < -1.0 || tap > 1.0 {
			t.Errorf("tap value %v out of reasonable range", tap)
		}
	}
}
