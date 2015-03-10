package pq

import (
	"strings"
	"testing"
)

func TestParseDateISO(t *testing.T) {
	for _, tt := range []struct {
		input            string
		year, month, day int
	}{
		{"2001-02-03", 2001, 2, 3},
		{"9999-99-99", 9999, 99, 99},
		{"0000001-02-03", 1, 2, 3},
		{"1010101-02-03", 1010101, 2, 3},
		{"0001-02-03 BC", 0, 2, 3},
		{"9999-99-99 BC", -9998, 99, 99},
	} {
		y, m, d, err := parseDateISO([]byte(tt.input))

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.input, err)
		}

		if y != tt.year {
			t.Errorf("Expected year to be %v for %q, got %v", tt.year, tt.input, y)
		}
		if m != tt.month {
			t.Errorf("Expected month to be %v for %q, got %v", tt.month, tt.input, m)
		}
		if d != tt.day {
			t.Errorf("Expected day to be %v for %q, got %v", tt.day, tt.input, d)
		}
	}
}

func TestParseDateISOError(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{"", "unexpected format"},
		{"1234567890", "unexpected format"},
		{"1234.67.90", "unexpected format"},
		{"1234-67.90", "unexpected format"},
		{"abcd-fg-ij", "expected number"},
		{"1234-fg-ij", "expected number"},
		{"1234-67-ij", "expected number"},
	} {
		_, _, _, err := parseDateISO([]byte(tt.input))

		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}

		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}
	}
}

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
