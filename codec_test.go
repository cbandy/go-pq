package pq

import (
	"strings"
	"testing"
	"time"
)

func TestAppendDateISO(t *testing.T) {
	for _, tt := range []struct {
		result           string
		year, month, day int
	}{
		{"0001-02-03", 1, 2, 3},
		{"2001-02-03", 2001, 2, 3},
		{"9999-99-99", 9999, 99, 99},
		{"1010101-02-03", 1010101, 2, 3},
		{"0001-02-03 BC", 0, 2, 3},
		{"9999-99-99 BC", -9998, 99, 99},
	} {
		result := appendDateISO([]byte{'x'}, tt.year, tt.month, tt.day)
		if string(result) != "x"+tt.result {
			t.Errorf("Expected %q to be appended for %v-%v-%v, got %q",
				tt.result, tt.year, tt.month, tt.day, result)
		}
	}
}

func TestAppendTime(t *testing.T) {
	for _, tt := range []struct {
		result                           string
		hour, minute, second, nanosecond int
	}{
		{"01:02:03", 1, 2, 3, 0},
		{"99:99:99", 99, 99, 99, 0},
		{"11:12:13.405060000", 11, 12, 13, 405060000},
		{"11:12:13.000004000", 11, 12, 13, 4000},
		{"11:12:13.000000004", 11, 12, 13, 4},
	} {
		result := appendTime([]byte{'x'}, tt.hour, tt.minute, tt.second, tt.nanosecond)
		if string(result) != "x"+tt.result {
			t.Errorf("Expected %q to be appended for %v:%v:%v.%v, got %q",
				tt.result, tt.hour, tt.minute, tt.second, tt.nanosecond, result)
		}
	}
}

func TestAppendTimestampISO(t *testing.T) {
	for _, tt := range []struct {
		result               string
		year, month, day     int
		hour, minute, second int
		nanosecond           int
	}{
		{"2001-02-03 04:05:06.007000000", 2001, 2, 3, 4, 5, 6, 7000000},
		{"9999-99-99 99:99:99", 9999, 99, 99, 99, 99, 99, 0},
		{"0001-02-03 04:05:06.007000000 BC", 0, 2, 3, 4, 5, 6, 7000000},
		{"9999-99-99 99:99:99 BC", -9998, 99, 99, 99, 99, 99, 0},
	} {
		result := appendTimestampISO([]byte{'x'},
			tt.year, tt.month, tt.day, tt.hour, tt.minute, tt.second, tt.nanosecond)

		if string(result) != "x"+tt.result {
			t.Errorf("Expected %q to be appended for %v-%v-%v %v:%v:%v.%v, got %q", tt.result,
				tt.year, tt.month, tt.day, tt.hour, tt.minute, tt.second, tt.nanosecond, result)
		}
	}
}

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
		{"1234-67-89 AD", "unexpected format"},
	} {
		_, _, _, err := parseDateISO([]byte(tt.input))

		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}

		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}

		if !strings.Contains(err.Error(), tt.input) {
			t.Errorf("Expected error to contain %q, got %q", tt.input, err)
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
		{"12:45:789", "unexpected format"},
		{"12:45:78.", "expected number"},
		{"12:45:78.xyz", "expected number"},
		{"12:45:78.0000000009", "unexpected format"},
	} {
		_, _, _, _, err := parseTime([]byte(tt.input))

		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}

		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}

		if !strings.Contains(err.Error(), tt.input) {
			t.Errorf("Expected error to contain %q, got %q", tt.input, err)
		}
	}
}

func TestParseTimestampISO(t *testing.T) {
	for _, tt := range []struct {
		input                string
		year, month, day     int
		hour, minute, second int
		nanosecond           int
	}{
		{"2001-02-03 04:05:06.007", 2001, 2, 3, 4, 5, 6, 7000000},
		{"9999-99-99 99:99:99", 9999, 99, 99, 99, 99, 99, 0},
		{"0001-02-03 04:05:06.007 BC", 0, 2, 3, 4, 5, 6, 7000000},
		{"9999-99-99 99:99:99 BC", -9998, 99, 99, 99, 99, 99, 0},
	} {
		y, mo, d, h, mi, s, ns, err := parseTimestampISO([]byte(tt.input))

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.input, err)
		}

		if y != tt.year {
			t.Errorf("Expected year to be %v for %q, got %v", tt.year, tt.input, y)
		}
		if mo != tt.month {
			t.Errorf("Expected month to be %v for %q, got %v", tt.month, tt.input, mo)
		}
		if d != tt.day {
			t.Errorf("Expected day to be %v for %q, got %v", tt.day, tt.input, d)
		}
		if h != tt.hour {
			t.Errorf("Expected hour to be %v for %q, got %v", tt.hour, tt.input, h)
		}
		if mi != tt.minute {
			t.Errorf("Expected minute to be %v for %q, got %v", tt.minute, tt.input, mi)
		}
		if s != tt.second {
			t.Errorf("Expected second to be %v for %q, got %v", tt.second, tt.input, s)
		}
		if ns != tt.nanosecond {
			t.Errorf("Expected nanosecond to be %v for %q, got %v", tt.nanosecond, tt.input, ns)
		}
	}
}

func TestParseTimestampISOError(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{"", "unexpected format"},
		{"2001-02-03", "unexpected format"},
		{"2001-02-03T04:05:06.007", "unexpected format"},
		{"abcd-fg-ij 04:05:06.007", "expected number"},
		{"2001-02-03 kl:mn:op.qrs", "expected number"},
	} {
		_, _, _, _, _, _, _, err := parseTimestampISO([]byte(tt.input))

		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}

		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}

		if !strings.Contains(err.Error(), tt.input) {
			t.Errorf("Expected error to contain %q, got %q", tt.input, err)
		}
	}
}

func TestParseTimestamptzISO(t *testing.T) {
	offset := func(h, m, s time.Duration) int {
		return int((h*time.Hour + m*time.Minute + s*time.Second) / time.Second)
	}
	for _, tt := range []struct {
		input                string
		year, month, day     int
		hour, minute, second int
		nanosecond, offset   int
	}{
		{"2001-02-03 04:05:06.007-08",
			2001, 2, 3, 4, 5, 6, 7000000, offset(-8, 0, 0)},
		{"2001-02-03 04:05:06.007-08:09",
			2001, 2, 3, 4, 5, 6, 7000000, offset(-8, -9, 0)},
		{"2001-02-03 04:05:06.007-08:09:10",
			2001, 2, 3, 4, 5, 6, 7000000, offset(-8, -9, -10)},
		{"9999-99-99 99:99:99+99:99:99",
			9999, 99, 99, 99, 99, 99, 0, offset(99, 99, 99)},

		{"0001-02-03 04:05:06.007-08 BC",
			0, 2, 3, 4, 5, 6, 7000000, offset(-8, 0, 0)},
		{"0001-02-03 04:05:06.007-08:09 BC",
			0, 2, 3, 4, 5, 6, 7000000, offset(-8, -9, 0)},
		{"0001-02-03 04:05:06.007-08:09:10 BC",
			0, 2, 3, 4, 5, 6, 7000000, offset(-8, -9, -10)},
		{"9999-99-99 99:99:99+99:99:99 BC",
			-9998, 99, 99, 99, 99, 99, 0, offset(99, 99, 99)},
	} {
		y, mo, d, h, mi, s, ns, o, err := parseTimestamptzISO([]byte(tt.input))

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.input, err)
		}

		if y != tt.year {
			t.Errorf("Expected year to be %v for %q, got %v", tt.year, tt.input, y)
		}
		if mo != tt.month {
			t.Errorf("Expected month to be %v for %q, got %v", tt.month, tt.input, mo)
		}
		if d != tt.day {
			t.Errorf("Expected day to be %v for %q, got %v", tt.day, tt.input, d)
		}
		if h != tt.hour {
			t.Errorf("Expected hour to be %v for %q, got %v", tt.hour, tt.input, h)
		}
		if mi != tt.minute {
			t.Errorf("Expected minute to be %v for %q, got %v", tt.minute, tt.input, mi)
		}
		if s != tt.second {
			t.Errorf("Expected second to be %v for %q, got %v", tt.second, tt.input, s)
		}
		if ns != tt.nanosecond {
			t.Errorf("Expected nanosecond to be %v for %q, got %v", tt.nanosecond, tt.input, ns)
		}
		if o != tt.offset {
			t.Errorf("Expected offset to be %v for %q, got %v", tt.offset, tt.input, o)
		}
	}
}

func TestParseTimestamptzISOError(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{"", "unexpected format"},
		{"2001-02-03", "unexpected format"},
		{"2001-02-03 04:05:06.007", "unexpected format"},
		{"2001-02-03 04:05:06.007Z", "expected number"},
		{"2001-02-03 04:05:06.007 08:09:10", "unexpected format"},
		{"2001-02-03T04:05:06.007+08:09:10", "unexpected format"},
		{"abcd-fg-ij 04:05:06.007+08:09:10", "expected number"},
		{"2001-02-03 kl:mn:op.qrs+08:09:10", "expected number"},
		{"2001-02-03 04:05:06.007+tu:vw:xy", "expected number"},
		{"2001-02-03 04:05:06.007+tu:vw:10", "expected number"},
		{"2001-02-03 04:05:06.007+tu:09:10", "expected number"},
	} {
		_, _, _, _, _, _, _, _, err := parseTimestamptzISO([]byte(tt.input))

		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}

		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}

		if !strings.Contains(err.Error(), tt.input) {
			t.Errorf("Expected error to contain %q, got %q", tt.input, err)
		}
	}
}
