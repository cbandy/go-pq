package pq

import (
	"bytes"
	"testing"
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
