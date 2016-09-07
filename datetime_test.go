package pq

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestDecodeTimestampIntegerError(t *testing.T) {
	_, _, err := decodeTimestampInteger([]byte{0x12, 0x34})

	if err == nil {
		t.Fatal("Expected error, got none")
	}
	if !strings.HasPrefix(err.Error(), "pq:") {
		t.Errorf("Expected error to start with %q, got %q", "pq:", err.Error())
	}
	if !strings.Contains(err.Error(), "bad length: 2") {
		t.Errorf("Expected error to contain length, got %q", err.Error())
	}
}

func TestDecodeTimestampBackend(t *testing.T) {
	db := openTestConn(t)
	defer db.Close()

	for _, tt := range []struct {
		s string
		t time.Time
	}{
		{"4300-07-08 09:10:11.012 BC", time.Date(-4299, time.July, 8, 9, 10, 11, 12000000, time.UTC)},
		{"0011-02-03 04:05:06.123 BC", time.Date(-10, time.February, 3, 4, 5, 6, 123000000, time.UTC)},
		{"2001-02-03 04:05:06.000001", time.Date(2001, time.February, 3, 4, 5, 6, 1000, time.UTC)},
		{"123456-02-03 04:05:06.0007", time.Date(123456, time.February, 3, 4, 5, 6, 700000, time.UTC)},
	} {
		var scanned interface{}
		err := db.QueryRow(`SELECT $1::timestamp`, tt.s).Scan(&scanned)
		if err != nil {
			t.Errorf("Expected no error for %q, got %v", tt.s, err)
			continue
		}
		if !reflect.DeepEqual(scanned, tt.t) {
			t.Errorf("Expected %v for %q, got %v", tt.t, tt.s, scanned)
		}
	}

	for _, tt := range []string{"infinity", "-infinity"} {
		var scanned interface{}
		err := db.QueryRow(`SELECT $1::timestamp`, tt).Scan(&scanned)
		if err != nil {
			t.Errorf("Expected no error for %q, got %v", tt, err)
			continue
		}
		if !reflect.DeepEqual(scanned, []byte(tt)) {
			t.Errorf("Expected []byte(%q), got %T(%q)", tt, scanned, scanned)
		}
	}
}
