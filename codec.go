package pq

import (
	"bytes"
	"fmt"
	"strconv"
)

// appendDateISO appends a date in the `yyyy{y...}-mm-dd[ BC]` format to the
// buffer and returns the extended buffer.
func appendDateISO(b []byte, year, month, day int) []byte {
	b, bc := appendDatePortionISO(b, year, month, day)
	return appendEraPortion(b, bc)
}

func appendDatePortionISO(b []byte, year, month, day int) (_ []byte, bc bool) {
	if year <= 0 {
		// Negate the year and add one.
		// See http://www.postgresql.org/docs/current/static/datetime-input-rules.html
		year = 1 - year
		bc = true
	}

	i := len(b)
	if year <= 9999 {
		b = append(b, "0000-00-00"...)
		writeDecimal(b[i:i+4], year)
		i += 4
	} else {
		b = strconv.AppendUint(b, uint64(year), 10)
		i = len(b)
		b = append(b, "-00-00"...)
	}

	writeDecimal(b[i+1:i+3], month)
	writeDecimal(b[i+4:i+6], day)

	return b, bc
}

func appendEraPortion(b []byte, bc bool) []byte {
	if bc {
		b = append(b, " BC"...)
	}

	return b
}

// appendTime appends a time in the `hh:mm:ss[.nnnnnnnnn]` format to the
// buffer and returns the extended buffer.
func appendTime(b []byte, hour, minute, second, nanosecond int) []byte {
	i := len(b)

	if nanosecond == 0 {
		b = append(b, "00:00:00"...)
	} else {
		b = append(b, "00:00:00.000000000"...)
		writeDecimal(b[i+9:i+18], nanosecond)
	}

	writeDecimal(b[i+0:i+2], hour)
	writeDecimal(b[i+3:i+5], minute)
	writeDecimal(b[i+6:i+8], second)

	return b
}

// appendTimestampISO appends a timestamp in the
// `yyyy{y...}-mm-dd hh:mm:ss[.nnnnnnnnn][ BC]` format to the buffer and
// returns the extended buffer.
func appendTimestampISO(b []byte, year, month, day, hour, minute, second, nanosecond int) []byte {
	b, bc := appendDatePortionISO(b, year, month, day)
	b = append(b, ' ')
	b = appendTime(b, hour, minute, second, nanosecond)
	return appendEraPortion(b, bc)
}

type parseErrorFunc func([]byte) error

func parseAtoI(src []byte, errFunc parseErrorFunc) (i int, err error) {
	for _, c := range src {
		i *= 10
		if '0' <= c && c <= '9' {
			i += int(c - '0')
		} else {
			err = errFunc(src)
		}
	}
	return
}

func parseDateGerman(src []byte) error {
	return fmt.Errorf("pq: unable to parse date; unexpected format for %q; not implemented", src)
}

// parseDateISO extracts the components of a date in the format
// `yyyy{y...}-mm-dd[ BC]`. Any other format results in an error.
func parseDateISO(src []byte) (year, month, day int, err error) {
	errAtoI := func(pos []byte) error {
		return fmt.Errorf("pq: unable to parse date; expected number at %q in %q", pos, src)
	}
	errFormat := func([]byte) error {
		return fmt.Errorf("pq: unable to parse date; unexpected format for %q", src)
	}

	year, month, day, unparsed, err := parseDatePortionISO(src, errAtoI, errFormat)
	if err != nil {
		return
	}

	if len(unparsed) > 0 {
		err = errFormat(src)
		return
	}

	return
}

// parseDatePortionISO extracts the date components of a value in the format
// `yyyy{y...}-mm-dd{...}[ BC]`. Any other format results in an error.
func parseDatePortionISO(src []byte, errAtoI, errFormat parseErrorFunc) (year, month, day int, remaining []byte, err error) {
	sepYearMonth := bytes.IndexByte(src, '-')
	if len(src) < 10 || sepYearMonth < 0 || src[sepYearMonth+3] != '-' {
		err = errFormat(src)
		return
	}

	if year, err = parseAtoI(src[0:sepYearMonth], errAtoI); err != nil {
		return
	}
	if month, err = parseAtoI(src[sepYearMonth+1:sepYearMonth+3], errAtoI); err != nil {
		return
	}
	if day, err = parseAtoI(src[sepYearMonth+4:sepYearMonth+6], errAtoI); err != nil {
		return
	}

	remaining = src[sepYearMonth+6:]

	// Dates before the current era are suffixed with " BC"
	if src[len(src)-3] == ' ' && src[len(src)-2] == 'B' && src[len(src)-1] == 'C' {
		// Negate the year and add one.
		// See http://www.postgresql.org/docs/current/static/datetime-input-rules.html
		year = 1 - year

		// Remove suffix
		remaining = remaining[:len(remaining)-3]
	}

	return
}

func parseDatePostgres(src []byte) error {
	return fmt.Errorf("pq: unable to parse date; ambiguous format for %q", src)
}

func parseDateSQL(src []byte) error {
	return fmt.Errorf("pq: unable to parse date; ambiguous format for %q", src)
}

// parseOffsetPortionISO extracts the offset component, in seconds, of a value
// in the format `{...}±hh[:mm[:ss]]`. Any other format results in an error.
func parseOffsetPortionISO(src []byte, errAtoI, errFormat parseErrorFunc) (offset int, remaining []byte, err error) {
	switch {
	case len(src) >= 9 && src[len(src)-6] == ':' && src[len(src)-3] == ':':
		var second int
		if second, err = parseAtoI(src[len(src)-2:], errAtoI); err != nil {
			return
		}
		offset += second
		src = src[:len(src)-3]
		fallthrough

	case len(src) >= 6 && src[len(src)-3] == ':':
		var minute int
		if minute, err = parseAtoI(src[len(src)-2:], errAtoI); err != nil {
			return
		}
		offset += 60 * minute
		src = src[:len(src)-3]
		fallthrough

	case len(src) >= 3:
		var hour int
		if hour, err = parseAtoI(src[len(src)-2:], errAtoI); err != nil {
			return
		}
		offset += 3600 * hour
		if src[len(src)-3] == '-' {
			offset = -offset
		} else if src[len(src)-3] != '+' {
			err = errFormat(src)
			return
		}
		src = src[:len(src)-3]

	default:
		err = errFormat(src)
		return
	}

	remaining = src
	return
}

// parseTime extracts the components of a time in the format
// `hh:mm:ss[.n{n...}]`. Any other format results in an error.
func parseTime(src []byte) (hour, minute, second, nanosecond int, err error) {
	errAtoI := func(pos []byte) error {
		return fmt.Errorf("pq: unable to parse time; expected number at %q in %q", pos, src)
	}
	errFormat := func([]byte) error {
		return fmt.Errorf("pq: unable to parse time; unexpected format for %q", src)
	}

	return parseTimePortion(src, errAtoI, errFormat)
}

// parseTimePortion extracts the time components of a value in the format
// `hh:mm:ss[.n{n...}]`. Any other format results in an error.
func parseTimePortion(src []byte, errAtoI, errFormat parseErrorFunc) (hour, minute, second, nanosecond int, err error) {
	if len(src) < 8 || src[2] != ':' || src[5] != ':' {
		err = errFormat(src)
		return
	}

	if hour, err = parseAtoI(src[0:2], errAtoI); err != nil {
		return
	}
	if minute, err = parseAtoI(src[3:5], errAtoI); err != nil {
		return
	}
	if second, err = parseAtoI(src[6:8], errAtoI); err != nil {
		return
	}

	if len(src) > 8 {
		if src[8] != '.' || len(src) > 18 {
			err = errFormat(src)
			return
		}

		if len(src) < 10 {
			err = errAtoI(src)
			return
		}

		if nanosecond, err = parseAtoI(src[9:], errAtoI); err != nil {
			return
		}

		// Scale to nanoseconds
		for i := len(src) - 8; i < 10; i++ {
			nanosecond *= 10
		}
	}

	return
}

// parseTimestampISO extracts the components of a timestamp in the format
// `yyyy{y...}-mm-dd hh:mm:ss[.n{n...}][ BC]`. Any other format results in
// an error.
func parseTimestampISO(src []byte) (year, month, day, hour, minute, second, nanosecond int, err error) {
	errAtoI := func(pos []byte) error {
		return fmt.Errorf("pq: unable to parse timestamp; expected number at %q in %q", pos, src)
	}
	errFormat := func([]byte) error {
		return fmt.Errorf("pq: unable to parse timestamp; unexpected format for %q", src)
	}

	year, month, day, unparsed, err := parseDatePortionISO(src, errAtoI, errFormat)
	if err != nil {
		return
	}

	if len(unparsed) < 1 || unparsed[0] != ' ' {
		err = errFormat(src)
		return
	}

	hour, minute, second, nanosecond, err = parseTimePortion(unparsed[1:], errAtoI, errFormat)
	return
}

func parseTimestampPostgres(src []byte) error {
	return fmt.Errorf("pq: unable to parse timestamp; unexpected format for %q; not implemented", src)
}

// parseTimestamptzISO extracts the components of a timestamptz in the format
// `yyyy{y...}-mm-dd hh:mm:ss[.n{n...}]±hh[:mm[:ss]][ BC]`. Any other format
// results in an error.
func parseTimestamptzISO(src []byte) (year, month, day, hour, minute, second, nanosecond, offset int, err error) {
	errAtoI := func(pos []byte) error {
		return fmt.Errorf("pq: unable to parse timestamptz; expected number at %q in %q", pos, src)
	}
	errFormat := func([]byte) error {
		return fmt.Errorf("pq: unable to parse timestamptz; unexpected format for %q", src)
	}

	year, month, day, unparsed, err := parseDatePortionISO(src, errAtoI, errFormat)
	if err != nil {
		return
	}

	offset, unparsed, err = parseOffsetPortionISO(unparsed, errAtoI, errFormat)
	if err != nil {
		return
	}

	if len(unparsed) < 1 || unparsed[0] != ' ' {
		err = errFormat(src)
		return
	}

	hour, minute, second, nanosecond, err = parseTimePortion(unparsed[1:], errAtoI, errFormat)
	return
}

// writeDecimal fills the destination buffer with the decimal representation of
// non-negative integer v. Values that do not fit are silently truncated.
func writeDecimal(dst []byte, v int) {
	for i := len(dst) - 1; i >= 0; i-- {
		q := v / 10
		dst[i] = byte('0' + (v - q*10))
		v = q
	}
}
