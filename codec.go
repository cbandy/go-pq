package pq

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
)

// Decode a byte slice from the bytea "escape" format.
func decodeByteaEscape(s []byte) (result []byte) {
	for len(s) > 0 {
		if s[0] == '\\' {
			// escaped '\\'
			if len(s) >= 2 && s[1] == '\\' {
				result = append(result, '\\')
				s = s[2:]
				continue
			}

			// '\\' followed by an octal number
			if len(s) < 4 {
				errorf("invalid bytea sequence %v", s)
			}
			u, err := strconv.ParseUint(string(s[1:4]), 8, 8)
			if err != nil {
				errorf("could not parse bytea value: %s", err)
			}
			result = append(result, byte(u))
			s = s[4:]
		} else {
			// We hit an unescaped, raw byte. Append as many as possible in one go.
			i := bytes.IndexByte(s, '\\')
			if i == -1 {
				result = append(result, s...)
				break
			}
			result = append(result, s[:i]...)
			s = s[i:]
		}
	}

	return
}

// Decode a byte slice from the bytea "hex" format.
func decodeByteaHex(s []byte) []byte {
	// Remove leading '\x'
	s = s[2:]
	result := make([]byte, hex.DecodedLen(len(s)))
	_, err := hex.Decode(result, s)
	if err != nil {
		errorf(err.Error())
	}
	return result
}

// Encode a byte slice to the bytea "escape" format.
func encodeByteaEscape(v []byte) (result []byte) {
	for _, b := range v {
		if b == '\\' {
			result = append(result, '\\', '\\')
		} else if b < 0x20 || b > 0x7e {
			result = append(result, fmt.Sprintf("\\%03o", b)...)
		} else {
			result = append(result, b)
		}
	}

	return result
}

// Encode a byte slice to the bytea "hex" format.
func encodeByteaHex(v []byte) []byte {
	result := make([]byte, 2+hex.EncodedLen(len(v)))
	result[0] = '\\'
	result[1] = 'x'
	hex.Encode(result[2:], v)
	return result
}
