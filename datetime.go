package pq

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"
)

// Clock represents a value of the PostgreSQL `time without time zone` type.
// It implements the sql.Scanner interface so it can be used as a scan
// destination.
type Clock struct {
	Hour, Minute, Second, Nanosecond int
}

// Scan implements the sql.Scanner interface.
func (c *Clock) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return c.scanBytes(src)

	case string:
		return c.scanBytes([]byte(src))

	case time.Time:
		return c.scanTime(src)
	}

	return fmt.Errorf("pq: cannot convert %T to Clock", src)
}

func (c *Clock) scanBytes(src []byte) (err error) {
	hour, min, sec, nsec, err := parseTime(src)

	if err == nil {
		*c = Clock{Hour: hour, Minute: min, Second: sec, Nanosecond: nsec}
	}

	return
}

func (c *Clock) scanTime(src time.Time) error {
	hour, min, sec := src.Clock()
	nsec := src.Nanosecond()

	*c = Clock{Hour: hour, Minute: min, Second: sec, Nanosecond: nsec}
	return nil
}

// Value implements the driver.Valuer interface.
func (c Clock) Value() (driver.Value, error) {
	return fmt.Sprintf("%02d:%02d:%02d.%09d", c.Hour, c.Minute, c.Second, c.Nanosecond), nil
}

// Date represents a value of the PostgreSQL `date` type. It implements the
// sql.Scanner interface so it can be used as a scan destination.
//
// A positive or negative value in Infinity represents the special value
// "infinity" or "-infinity", respectively.
type Date struct {
	Infinity int
	Year     int
	Month    time.Month
	Day      int
}

// Scan implements the sql.Scanner interface.
func (d *Date) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return d.scanBytes(src)

	case string:
		return d.scanBytes([]byte(src))

	case time.Time:
		return d.scanTime(src)
	}

	return fmt.Errorf("pq: cannot convert %T to Date", src)
}

func (d *Date) scanBytes(src []byte) (err error) {
	if len(src) > 2 {
		switch {
		case 'n' == src[2] && bytes.Equal(src, []byte{'-', 'i', 'n', 'f', 'i', 'n', 'i', 't', 'y'}):
			*d = Date{Infinity: -1}
			return

		case 'f' == src[2] && bytes.Equal(src, []byte{'i', 'n', 'f', 'i', 'n', 'i', 't', 'y'}):
			*d = Date{Infinity: 1}
			return

		case '0' <= src[2] && src[2] <= '9':
			year, month, day, err := parseDateISO(src)
			if err == nil {
				*d = Date{Year: year, Month: time.Month(month), Day: day}
			}
			return err

		case '.' == src[2]:
			return parseDateGerman(src)

		case '/' == src[2]:
			return parseDateSQL(src)
		}
	}

	return parseDatePostgres(src)
}

func (d *Date) scanTime(src time.Time) error {
	year, month, day := src.Date()
	*d = Date{Year: year, Month: month, Day: day}
	return nil
}

// Value implements the driver.Valuer interface.
func (d Date) Value() (driver.Value, error) {
	switch {
	case d.Infinity < 0:
		return "-infinity", nil

	case d.Infinity > 0:
		return "infinity", nil

	default:
		return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day), nil
	}
}

// Timestamp represents a value of the PostgreSQL `timestamp without time zone`
// type. It supports the special values "infinity" and "-infinity" and
// implements the sql.Scanner interface so it can be used as a scan destination.
type Timestamp struct {
	Date
	Clock
}

// Scan implements the sql.Scanner interface.
func (t *Timestamp) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return t.scanBytes(src)

	case string:
		return t.scanBytes([]byte(src))

	case time.Time:
		return t.scanTime(src)
	}

	return fmt.Errorf("pq: cannot convert %T to Timestamp", src)
}

func (t *Timestamp) scanBytes(src []byte) (err error) {
	if len(src) > 2 {
		switch {
		case 'n' == src[2] && bytes.Equal(src, []byte{'-', 'i', 'n', 'f', 'i', 'n', 'i', 't', 'y'}):
			t.Date = Date{Infinity: -1}
			t.Clock = Clock{}
			return

		case 'f' == src[2] && bytes.Equal(src, []byte{'i', 'n', 'f', 'i', 'n', 'i', 't', 'y'}):
			t.Date = Date{Infinity: 1}
			t.Clock = Clock{}
			return

		case '0' <= src[2] && src[2] <= '9':
			year, month, day, hour, min, sec, nsec, err := parseTimestampISO(src)
			if err == nil {
				t.Date = Date{Year: year, Month: time.Month(month), Day: day}
				t.Clock = Clock{Hour: hour, Minute: min, Second: sec, Nanosecond: nsec}
			}
			return err

		case '.' == src[2]:
			return parseDateGerman(src)

		case '/' == src[2]:
			return parseDateSQL(src)
		}
	}

	return parseTimestampPostgres(src)
}

func (t *Timestamp) scanTime(src time.Time) error {
	t.Date.scanTime(src)
	t.Clock.scanTime(src)
	return nil
}

// TimestampTZ represents a value of the PostgreSQL `timestamp with time zone`
// type. It implements the sql.Scanner interface so it can be used as a scan
// destination.
//
// A positive or negative value in Infinity represents the special value
// "infinity" or "-infinity", respectively.
type TimestampTZ struct {
	Infinity int
	Time     time.Time
}
