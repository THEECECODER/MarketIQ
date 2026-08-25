package marketdata

import "testing"

func TestParsePrice(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"78914.86", 7891486},
		{"100", 10000},
		{"100.5", 10050},
		{"0.01", 1},
	}

	for _, tt := range tests {
		got, err := parsePrice(tt.input)

		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tt.input, err)
		}

		if got != tt.expected {
			t.Fatalf(
				"parsePrice(%q): expected %d, got %d",
				tt.input,
				tt.expected,
				got,
			)
		}
	}
}
func TestParseQuantity(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"0.7402", 74020},
		{"1", 100000},
		{"1.5", 150000},
		{"0.0001", 10},
		{"0.00014", 14},
		{"0.04497", 4497},
	}

	for _, tt := range tests {
		got, err := parseQuantity(tt.input)

		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tt.input, err)
		}

		if got != tt.expected {
			t.Fatalf(
				"parseQuantity(%q): expected %d, got %d",
				tt.input,
				tt.expected,
				got,
			)
		}
	}
}
