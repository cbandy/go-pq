package pq

import (
	"fmt"
)

func parseAtoI(src []byte, errFunc func([]byte) error) (i int, err error) {
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
