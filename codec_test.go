package pq

import (
	"strings"
	"testing"
)

func TestParseTime(t *testing.T) {
	for _, tt := range []struct {
		input                            string
		hour, minute, second, nanosecond int
	}{
		{"01:02:03", 1, 2, 3, 0},
		{"99:99:99", 99, 99, 99, 0},
		{"11:12:13.40506", 11, 12, 13, 405060000},
		{"11:12:13.000004", 11, 12, 13, 4000},
		{"11:12:13.000000004", 11, 12, 13, 4},
	} {
		h, m, s, ns, err := parseTime([]byte(tt.input))

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.input, err)
		}

		if h != tt.hour {
			t.Errorf("Expected hour to be %v for %q, got %v", tt.hour, tt.input, h)
		}
		if m != tt.minute {
			t.Errorf("Expected minute to be %v for %q, got %v", tt.minute, tt.input, m)
		}
		if s != tt.second {
			t.Errorf("Expected second to be %v for %q, got %v", tt.second, tt.input, s)
		}
		if ns != tt.nanosecond {
			t.Errorf("Expected nanosecond to be %v for %q, got %v", tt.nanosecond, tt.input, ns)
		}
	}
}

func TestParseTimeError(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{"", "unexpected format"},
		{"12345678", "unexpected format"},
		{"12.45.78", "unexpected format"},
		{"12:45.78", "unexpected format"},
		{"ab:de:gh", "expected number"},
		{"12:de:gh", "expected number"},
		{"12:45:gh", "expected number"},
		{"12:45:789", "expected '.'"},
		{"12:45:78.", "expected number"},
		{"12:45:78.xyz", "expected number"},
	} {
		_, _, _, _, err := parseTime([]byte(tt.input))

		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}

		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}
	}
}
