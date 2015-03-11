package pq

import (
	"bytes"
	"fmt"
)

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
	errAtoI := func(src []byte) error {
		return fmt.Errorf("pq: unable to parse date; expected number at %q", src)
	}
	errFormat := func(src []byte) error {
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

// parseTime extracts the components of a time in the format
// `hh:mm:ss[.n{n...}]`. Any other format results in an error.
func parseTime(src []byte) (hour, minute, second, nanosecond int, err error) {
	errAtoI := func(src []byte) error {
		return fmt.Errorf("pq: unable to parse time; expected number at %q", src)
	}

	if len(src) < 8 || src[2] != ':' || src[5] != ':' {
		err = fmt.Errorf("pq: unable to parse time; unexpected format for %q", src)
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
		if src[8] != '.' {
			err = fmt.Errorf("pq: unable to parse time; expected '.' in %q", src)
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
	errAtoI := func(src []byte) error {
		return fmt.Errorf("pq: unable to parse timestamp; expected number at %q", src)
	}
	errFormat := func(src []byte) error {
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

	hour, minute, second, nanosecond, err = parseTime(unparsed[1:])
	return
}

func parseTimestampPostgres(src []byte) error {
	return fmt.Errorf("pq: unable to parse timestamp; unexpected format for %q; not implemented", src)
}
