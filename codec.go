package pq

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
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

// Decode a Time from the "ISO" format.
func decodeTimestamptzISO(src []byte) time.Time {
	atoi := func(s []byte) (result int) {
		for i := 0; i < len(s); i++ {
			result = result*10 + int(s[i]-'0')
		}
		return
	}

	sepYearMonth := bytes.IndexByte(src, '-')
	year := atoi(src[:sepYearMonth])
	src = src[sepYearMonth:]

	// Skips a separator and converts two digits
	nextTwoDigits := func() (result int) {
		result = atoi(src[1:3])
		src = src[3:]
		return
	}

	month := nextTwoDigits()
	day := nextTwoDigits()
	hour := nextTwoDigits()
	minute := nextTwoDigits()
	second := nextTwoDigits()
	nanosecond, offset := 0, 0

	// Time before current era is suffixed with BC
	if src[len(src)-1] == 'C' {
		// Negate the year and add one.
		// See http://www.postgresql.org/docs/current/static/datetime-input-rules.html
		year = 1 - year

		// Strip " BC"
		src = src[:len(src)-3]
	}

	// Offset from UTC is formatted ±hh[:mm[:ss]]
	switch {
	case len(src) > 6 && src[len(src)-6] == ':':
		offset += atoi(src[len(src)-2:])
		src = src[:len(src)-3]
		fallthrough

	case len(src) > 3 && src[len(src)-3] == ':':
		offset += 60 * atoi(src[len(src)-2:])
		src = src[:len(src)-3]
		fallthrough

	default:
		offset += 3600 * atoi(src[len(src)-2:])
		if src[len(src)-3] == '-' {
			offset = -offset
		}
		src = src[:len(src)-3]
	}

	// Fractional seconds
	if len(src) > 1 {
		// Skip fraction separator
		i := 1
		for ; i < len(src); i++ {
			nanosecond = nanosecond*10 + int(src[i]-'0')
		}
		// Scale to nanosecnds
		for ; i < 10; i++ {
			nanosecond *= 10
		}
	}

	return time.Date(
		year, time.Month(month), day,
		hour, minute, second, nanosecond,
		globalLocationCache.getLocation(offset))
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
