package common

import (
	"testing"
)

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		name      string
		n         float64
		precision int
		expected  float64
	}{
		{"precision 2 round down", 3.14159, 2, 3.14},
		{"precision 3 round up", 3.14159, 3, 3.142},
		{"precision 0", 3.14159, 0, 3},
		{"precision 1", 3.75, 1, 3.8},
		{"negative number", -2.6789, 2, -2.68},
		{"zero value", 0.0, 2, 0.0},
		{"large precision", 1.123456789, 8, 1.12345679},
		{"whole number", 5.0, 2, 5.0},
		{"small number", 0.001234, 4, 0.0012},
		{"precision 5 no rounding needed", 1.5, 5, 1.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFloat(tt.n, tt.precision)
			t.Log(result)
			if result != tt.expected {
				t.Errorf("FormatFloat(%v, %d) = %v, want %v", tt.n, tt.precision, result, tt.expected)
			}
		})
	}
}
