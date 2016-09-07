package pq

import (
	"encoding/binary"
	"fmt"
	"time"
)

const (
	pgToUnix int64 = 946684800 * 1000000 // microseconds
	unixToPg int64 = -pgToUnix
)

func decodeTimestampFloat([]byte) (int, time.Time, error) {
	return 0, time.Time{}, fmt.Errorf("pq: unable to decode timestamp; float not implemented")
}

// decodeTimestampInteger interprets the binary format of a timestamp or timestamptz
// when the backend stores them as eight-byte integers.
func decodeTimestampInteger(src []byte) (int, time.Time, error) {
	if len(src) != 8 {
		return 0, time.Time{}, fmt.Errorf("pq: unable to decode timestamp; bad length: %d", len(src))
	}

	usec := int64(binary.BigEndian.Uint64(src))

	if usec <= -9223372036854775808 {
		return -1, time.Time{}, nil
	}
	if usec >= 9223372036854775807 {
		return 1, time.Time{}, nil
	}

	usec += pgToUnix

	return 0, time.Unix(usec/1000000, (usec%1000000)*1000).In(time.UTC), nil
}

func encodeTimestampFloat(time.Time) ([]byte, error) {
	return nil, fmt.Errorf("pq: unable to encode timestamp; float not implemented")
}

func encodeTimestampInteger(time.Time) ([]byte, error) {
	return nil, fmt.Errorf("pq: unable to encode timestamp; nanoseconds would be lost")
}
