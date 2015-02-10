package pq

import (
	"bytes"
	"testing"
	"time"
)

type byteaTest struct{ raw, encoded []byte }

var byteaEscapeTests []byteaTest = []byteaTest{
	{[]byte{}, []byte(``)},
	{[]byte{0x0}, []byte(`\000`)},
	{[]byte{0xde, 0xad, 0xbe, 0xef}, []byte(`\336\255\276\357`)},
	{[]byte{'a', 's', 0x0, 'c', 0x0, 'i', 'i'}, []byte(`as\000c\000ii`)},
}

func TestDecodeByteaEscape(t *testing.T) {
	for _, tt := range byteaEscapeTests {
		result := decodeByteaEscape(tt.encoded)
		if !bytes.Equal(result, tt.raw) {
			t.Errorf("Expected %x, got %x", tt.raw, result)
		}
	}
}

func TestEncodeByteaEscape(t *testing.T) {
	for _, tt := range byteaEscapeTests {
		result := encodeByteaEscape(tt.raw)
		if !bytes.Equal(result, tt.encoded) {
			t.Errorf("Expected %x, got %x", tt.encoded, result)
		}
	}
}

var byteaHexTests []byteaTest = []byteaTest{
	{[]byte{}, []byte(`\x`)},
	{[]byte{0x0}, []byte(`\x00`)},
	{[]byte{0xde, 0xad, 0xbe, 0xef}, []byte(`\xdeadbeef`)},
}

func TestDecodeByteaHex(t *testing.T) {
	for _, tt := range byteaHexTests {
		result := decodeByteaHex(tt.encoded)
		if !bytes.Equal(result, tt.raw) {
			t.Errorf("Expected %x, got %x", tt.raw, result)
		}
	}
}

func TestEncodeByteaHex(t *testing.T) {
	for _, tt := range byteaHexTests {
		result := encodeByteaHex(tt.raw)
		if !bytes.Equal(result, tt.encoded) {
			t.Errorf("Expected %x, got %x", tt.encoded, result)
		}
	}
}

type timeTest struct {
	raw         time.Time
	fromBackend []byte
	toBackend   []byte
}

var timestamptzISOTests = []timeTest{
	{time.Date(1, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", 0)),
		[]byte(`0001-02-03 04:05:06.123456789+00`),
		[]byte(`0001-02-03T04:05:06.123456789Z`)},
	{time.Date(1, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", 2*60*60)),
		[]byte(`0001-02-03 04:05:06.123456789+02`),
		[]byte(`0001-02-03T04:05:06.123456789+02:00`)},
	{time.Date(1, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", -6*60*60)),
		[]byte(`0001-02-03 04:05:06.123456789-06`),
		[]byte(`0001-02-03T04:05:06.123456789-06:00`)},
	{time.Date(1, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", 7*60*60+30*60+9)),
		[]byte(`0001-02-03 04:05:06.123456789+07:30:09`),
		[]byte(`0001-02-03T04:05:06.123456789+07:30:09`)},

	{time.Date(1, time.February, 3, 4, 5, 6, 0, time.FixedZone("", 0)),
		[]byte(`0001-02-03 04:05:06+00`),
		[]byte(`0001-02-03T04:05:06Z`)},
	{time.Date(1, time.February, 3, 4, 5, 6, 1000, time.FixedZone("", 0)),
		[]byte(`0001-02-03 04:05:06.000001+00`),
		[]byte(`0001-02-03T04:05:06.000001Z`)},
	{time.Date(1, time.February, 3, 4, 5, 6, 1000000, time.FixedZone("", 0)),
		[]byte(`0001-02-03 04:05:06.001+00`),
		[]byte(`0001-02-03T04:05:06.001Z`)},

	{time.Date(10000, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", 0)),
		[]byte(`10000-02-03 04:05:06.123456789+00`),
		[]byte(`10000-02-03T04:05:06.123456789Z`)},
	{time.Date(10000, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", 2*60*60)),
		[]byte(`10000-02-03 04:05:06.123456789+02`),
		[]byte(`10000-02-03T04:05:06.123456789+02:00`)},
	{time.Date(10000, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", -6*60*60)),
		[]byte(`10000-02-03 04:05:06.123456789-06`),
		[]byte(`10000-02-03T04:05:06.123456789-06:00`)},
	{time.Date(10000, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", 7*60*60+30*60+9)),
		[]byte(`10000-02-03 04:05:06.123456789+07:30:09`),
		[]byte(`10000-02-03T04:05:06.123456789+07:30:09`)},

	{time.Date(10000, time.February, 3, 4, 5, 6, 0, time.FixedZone("", 0)),
		[]byte(`10000-02-03 04:05:06+00`),
		[]byte(`10000-02-03T04:05:06Z`)},
	{time.Date(10000, time.February, 3, 4, 5, 6, 1000, time.FixedZone("", 0)),
		[]byte(`10000-02-03 04:05:06.000001+00`),
		[]byte(`10000-02-03T04:05:06.000001Z`)},
	{time.Date(10000, time.February, 3, 4, 5, 6, 1000000, time.FixedZone("", 0)),
		[]byte(`10000-02-03 04:05:06.001+00`),
		[]byte(`10000-02-03T04:05:06.001Z`)},

	{time.Date(0, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", 0)),
		[]byte(`0001-02-03 04:05:06.123456789+00 BC`),
		[]byte(`0001-02-03T04:05:06.123456789Z BC`)},
	{time.Date(0, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", 2*60*60)),
		[]byte(`0001-02-03 04:05:06.123456789+02 BC`),
		[]byte(`0001-02-03T04:05:06.123456789+02:00 BC`)},
	{time.Date(0, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", -6*60*60)),
		[]byte(`0001-02-03 04:05:06.123456789-06 BC`),
		[]byte(`0001-02-03T04:05:06.123456789-06:00 BC`)},
	{time.Date(0, time.February, 3, 4, 5, 6, 123456789, time.FixedZone("", 7*60*60+30*60+9)),
		[]byte(`0001-02-03 04:05:06.123456789+07:30:09 BC`),
		[]byte(`0001-02-03T04:05:06.123456789+07:30:09 BC`)},

	{time.Date(0, time.February, 3, 4, 5, 6, 0, time.FixedZone("", 0)),
		[]byte(`0001-02-03 04:05:06+00 BC`),
		[]byte(`0001-02-03T04:05:06Z BC`)},
	{time.Date(0, time.February, 3, 4, 5, 6, 1000, time.FixedZone("", 0)),
		[]byte(`0001-02-03 04:05:06.000001+00 BC`),
		[]byte(`0001-02-03T04:05:06.000001Z BC`)},
	{time.Date(0, time.February, 3, 4, 5, 6, 1000000, time.FixedZone("", 0)),
		[]byte(`0001-02-03 04:05:06.001+00 BC`),
		[]byte(`0001-02-03T04:05:06.001Z BC`)},
}

func TestDecodeTimestamptzISO(t *testing.T) {
	for _, tt := range timestamptzISOTests {
		result := decodeTimestamptzISO(tt.fromBackend)
		if !tt.raw.Equal(result) || tt.raw.Format("-0700 MST") != result.Format("-0700 MST") {
			t.Errorf("Expected %v, got %v", tt.raw, result)
		}
	}
}
