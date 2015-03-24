package pq

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClockScanUnsupported(t *testing.T) {
	var clock Clock
	err := clock.Scan(true)

	if err == nil {
		t.Fatal("Expected error when scanning from bool")
	}
	if !strings.Contains(err.Error(), "bool to Clock") {
		t.Errorf("Expected type to be mentioned when scanning, got %q", err)
	}
}

func TestClockScanTime(t *testing.T) {
	clock := Clock{9, 9, 9, 9}
	err := clock.Scan(time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if clock != (Clock{Hour: 4, Minute: 5, Second: 6, Nanosecond: 7}) {
		t.Errorf("Expected 04:05:06.000000007, got %+v", clock)
	}
}

func BenchmarkClockScanTime(b *testing.B) {
	var clock Clock
	var x, _ interface{} = time.Parse("15:04:05", `01:02:03`)
	var y, _ interface{} = time.Parse("15:04:05", `01:02:03.004005`)

	for i := 0; i < b.N; i++ {
		clock.Scan(x)
		clock.Scan(y)
	}
}

var ClockStringTests = []struct {
	str   string
	clock Clock
}{
	{`04:05:06`, Clock{Hour: 4, Minute: 5, Second: 6}},
	{`04:05:06.007`, Clock{Hour: 4, Minute: 5, Second: 6, Nanosecond: 7000000}},
	{`04:05:06.000007`, Clock{Hour: 4, Minute: 5, Second: 6, Nanosecond: 7000}},
}

func TestClockScanBytes(t *testing.T) {
	for _, tt := range ClockStringTests {
		bytes := []byte(tt.str)
		clock := Clock{9, 9, 9, 9}
		err := clock.Scan(bytes)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", bytes, err)
		}
		if clock != tt.clock {
			t.Errorf("Expected %+v for %q, got %+v", tt.clock, bytes, clock)
		}
	}
}

func BenchmarkClockScanBytes(b *testing.B) {
	var clock Clock
	var x interface{} = []byte(`01:02:03`)
	var y interface{} = []byte(`01:02:03.004005`)

	for i := 0; i < b.N; i++ {
		clock.Scan(x)
		clock.Scan(y)
	}
}

func TestClockScanString(t *testing.T) {
	for _, tt := range ClockStringTests {
		clock := Clock{9, 9, 9, 9}
		err := clock.Scan(tt.str)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.str, err)
		}
		if clock != tt.clock {
			t.Errorf("Expected %+v, got %+v", tt.clock, clock)
		}
	}
}

func TestClockScanError(t *testing.T) {
	clock := Clock{9, 9, 9, 9}
	err := clock.Scan("")

	if err == nil {
		t.Error("Expected error, got none")
	}
	if clock != (Clock{9, 9, 9, 9}) {
		t.Errorf("Expected destination not to change, got %+v", clock)
	}
}

func TestClockValue(t *testing.T) {
	for _, tt := range []struct {
		str   string
		clock Clock
	}{
		{`04:05:06`, Clock{Hour: 4, Minute: 5, Second: 6}},
		{`04:05:06.007000000`, Clock{Hour: 4, Minute: 5, Second: 6, Nanosecond: 7000000}},
		{`04:05:06.000007000`, Clock{Hour: 4, Minute: 5, Second: 6, Nanosecond: 7000}},
		{`04:05:06.000000007`, Clock{Hour: 4, Minute: 5, Second: 6, Nanosecond: 7}},
	} {
		value, err := tt.clock.Value()

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.clock, err)
		}
		if string(value.([]byte)) != tt.str {
			t.Errorf("Expected %v, got %v", tt.str, value)
		}
	}
}

func BenchmarkClockValue(b *testing.B) {
	x := Clock{Hour: 4, Minute: 5, Second: 6}
	y := Clock{Hour: 4, Minute: 5, Second: 6, Nanosecond: 7000}

	for i := 0; i < b.N; i++ {
		x.Value()
		y.Value()
	}
}

func TestDateScanUnsupportedType(t *testing.T) {
	var date Date
	err := date.Scan(true)

	if err == nil {
		t.Fatal("Expected error when scanning from bool")
	}
	if !strings.Contains(err.Error(), "bool to Date") {
		t.Errorf("Expected type to be mentioned when scanning, got %q", err)
	}
}

func TestDateScanUnsupportedFormat(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{`02/03/2001`, "ambiguous format"}, // SQL, MDY
		{`03/02/2001`, "ambiguous format"}, // SQL, DMY
		{`02-03-2001`, "ambiguous format"}, // Postgres, MDY
		{`03-02-2001`, "ambiguous format"}, // Postgres, DMY
		{`03.02.2001`, "not implemented"},  // German
	} {
		date := Date{9, 9, 9, 9}
		err := date.Scan(tt.input)

		if err == nil {
			t.Fatal("Expected error, got none")
		}

		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}
	}
}

func TestDateScanTime(t *testing.T) {
	date := Date{9, 9, 9, 9}
	err := date.Scan(time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if date != (Date{Year: 2001, Month: 2, Day: 3}) {
		t.Errorf("Expected 2001-02-03, got %v", date)
	}
}

func BenchmarkDateScanTime(b *testing.B) {
	var date Date
	var x, _ interface{} = time.Parse("2006-01-02", `2001-02-03`)

	for i := 0; i < b.N; i++ {
		date.Scan(x)
	}
}

var DateStringTests = []struct {
	str  string
	date Date
}{
	{`infinity`, Date{Infinity: 1}},
	{`-infinity`, Date{Infinity: -1}},
	{`2001-02-03`, Date{Year: 2001, Month: 2, Day: 3}},
	{`4000-05-06 BC`, Date{Year: -3999, Month: 5, Day: 6}},
}

func TestDateScanBytes(t *testing.T) {
	for _, tt := range DateStringTests {
		bytes := []byte(tt.str)
		date := Date{9, 9, 9, 9}
		err := date.Scan(bytes)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", bytes, err)
		}
		if date != tt.date {
			t.Errorf("Expected %+v for %q, got %+v", tt.date, bytes, date)
		}
	}
}

func BenchmarkDateScanBytesISO(b *testing.B) {
	var date Date
	var x interface{} = []byte(`2001-02-03`)
	var y interface{} = []byte(`2001-02-03 BC`)

	for i := 0; i < b.N; i++ {
		date.Scan(x)
		date.Scan(y)
	}
}

func BenchmarkDateScanBytesInfinity(b *testing.B) {
	var date Date
	var x interface{} = []byte(`-infinity`)
	var y interface{} = []byte(`infinity`)

	for i := 0; i < b.N; i++ {
		date.Scan(x)
		date.Scan(y)
	}
}

func TestDateScanString(t *testing.T) {
	for _, tt := range DateStringTests {
		date := Date{9, 9, 9, 9}
		err := date.Scan(tt.str)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.str, err)
		}
		if date != tt.date {
			t.Errorf("Expected %+v, got %+v", tt.date, date)
		}
	}
}

func TestDateScanError(t *testing.T) {
	date := Date{9, 9, 9, 9}
	err := date.Scan("")

	if err == nil {
		t.Error("Expected error, got none")
	}
	if date != (Date{9, 9, 9, 9}) {
		t.Errorf("Expected destination not to change, got %+v", date)
	}
}

func TestDateValue(t *testing.T) {
	for _, tt := range DateStringTests {
		value, err := tt.date.Value()

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.date, err)
		}
		switch v := value.(type) {
		case string:
			if value != tt.str {
				t.Errorf("Expected %v, got %v", tt.str, value)
			}
		case []byte:
			if value = string(v); value != tt.str {
				t.Errorf("Expected %v, got %v", tt.str, value)
			}
		}
	}
}

func BenchmarkDateValue(b *testing.B) {
	x := Date{Year: 2001, Month: 2, Day: 3}
	y := Date{Year: 10000, Month: 2, Day: 3}
	z := Date{Year: -4000, Month: 2, Day: 3}

	for i := 0; i < b.N; i++ {
		x.Value()
		y.Value()
		z.Value()
	}
}

func BenchmarkDateValueInfinity(b *testing.B) {
	x := Date{Infinity: -1}
	y := Date{Infinity: 1}

	for i := 0; i < b.N; i++ {
		x.Value()
		y.Value()
	}
}

func TestTimestampScanUnsupportedType(t *testing.T) {
	var ts Timestamp
	err := ts.Scan(true)

	if err == nil {
		t.Fatal("Expected error when scanning from bool")
	}
	if !strings.Contains(err.Error(), "bool to Timestamp") {
		t.Errorf("Expected type to be mentioned when scanning, got %q", err)
	}
}

func TestTimestampScanUnsupportedFormat(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{`02/03/2001 04:05:06.007`, "ambiguous format"},     // SQL, MDY
		{`03/02/2001 04:05:06.007`, "ambiguous format"},     // SQL, DMY
		{`Sat Feb 03 04:05:06.007 2001`, "not implemented"}, // Postgres, MDY
		{`Sat 03 Feb 04:05:06.007 2001`, "not implemented"}, // Postgres, DMY
		{`03.02.2001 04:05:06.007`, "not implemented"},      // German
	} {
		ts := Timestamp{Date{9, 9, 9, 9}, Clock{9, 9, 9, 9}}
		err := ts.Scan(tt.input)

		if err == nil {
			t.Fatal("Expected error for %q, got none", tt.input)
		}

		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}
	}
}

func TestTimestampScanTime(t *testing.T) {
	ts := Timestamp{Date{9, 9, 9, 9}, Clock{9, 9, 9, 9}}
	err := ts.Scan(time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if ts != (Timestamp{
		Date{Year: 2001, Month: time.February, Day: 3},
		Clock{Hour: 4, Minute: 5, Second: 6, Nanosecond: 7},
	}) {
		t.Errorf("Expected 2001-02-03 04:05:06.000000007, got %+v", ts)
	}
}

func BenchmarkTimestampScanTime(b *testing.B) {
	var ts Timestamp
	var x, _ interface{} = time.Parse("2006-01-02 15:04:05", `2001-02-03 04:05:06`)
	var y, _ interface{} = time.Parse("2006-01-02 15:04:05", `2001-02-03 04:05:06.007008`)

	for i := 0; i < b.N; i++ {
		ts.Scan(x)
		ts.Scan(y)
	}
}

var TimestampStringTests = []struct {
	str       string
	timestamp Timestamp
}{
	{`infinity`, Timestamp{Date: Date{Infinity: 1}}},
	{`-infinity`, Timestamp{Date: Date{Infinity: -1}}},
	{`2001-02-03 04:05:06`,
		Timestamp{
			Date:  Date{Year: 2001, Month: 2, Day: 3},
			Clock: Clock{Hour: 4, Minute: 5, Second: 6}}},
	{`2001-02-03 04:05:06.007`,
		Timestamp{
			Date:  Date{Year: 2001, Month: 2, Day: 3},
			Clock: Clock{Hour: 4, Minute: 5, Second: 6, Nanosecond: 7000000}}},
	{`2001-02-03 04:05:06 BC`,
		Timestamp{
			Date:  Date{Year: -2000, Month: 2, Day: 3},
			Clock: Clock{Hour: 4, Minute: 5, Second: 6}}},
	{`2001-02-03 04:05:06.007 BC`,
		Timestamp{
			Date:  Date{Year: -2000, Month: 2, Day: 3},
			Clock: Clock{Hour: 4, Minute: 5, Second: 6, Nanosecond: 7000000}}},
}

func TestTimestampScanBytes(t *testing.T) {
	for _, tt := range TimestampStringTests {
		bytes := []byte(tt.str)
		ts := Timestamp{Date{9, 9, 9, 9}, Clock{9, 9, 9, 9}}
		err := ts.Scan(bytes)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", bytes, err)
		}
		if ts != tt.timestamp {
			t.Errorf("Expected %+v for %q, got %+v", tt.timestamp, bytes, ts)
		}
	}
}

func BenchmarkTimestampScanBytesISO(b *testing.B) {
	var ts Timestamp
	var w interface{} = []byte(`2001-02-03 04:05:06`)
	var x interface{} = []byte(`2001-02-03 04:05:06.007008`)
	var y interface{} = []byte(`2001-02-03 04:05:06 BC`)
	var z interface{} = []byte(`2001-02-03 04:05:06.007008 BC`)

	for i := 0; i < b.N; i++ {
		ts.Scan(w)
		ts.Scan(x)
		ts.Scan(y)
		ts.Scan(z)
	}
}

func BenchmarkTimestampScanBytesInfinity(b *testing.B) {
	var ts Timestamp
	var x interface{} = []byte(`-infinity`)
	var y interface{} = []byte(`infinity`)

	for i := 0; i < b.N; i++ {
		ts.Scan(x)
		ts.Scan(y)
	}
}

func TestTimestampScanString(t *testing.T) {
	for _, tt := range TimestampStringTests {
		ts := Timestamp{Date{9, 9, 9, 9}, Clock{9, 9, 9, 9}}
		err := ts.Scan(tt.str)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.str, err)
		}
		if ts != tt.timestamp {
			t.Errorf("Expected %+v for %q, got %+v", tt.timestamp, tt.str, ts)
		}
	}
}

func TestTimestampScanError(t *testing.T) {
	ts := Timestamp{Date{9, 9, 9, 9}, Clock{9, 9, 9, 9}}
	err := ts.Scan("")

	if err == nil {
		t.Error("Expected error, got none")
	}
	if ts != (Timestamp{Date{9, 9, 9, 9}, Clock{9, 9, 9, 9}}) {
		t.Errorf("Expected destination not to change, got %+v", ts)
	}
}

func TestTimestampTZScanUnsupportedType(t *testing.T) {
	var tstz TimestampTZ
	err := tstz.Scan(true)

	if err == nil {
		t.Fatal("Expected error when scanning from bool")
	}
	if !strings.Contains(err.Error(), "bool to TimestampTZ") {
		t.Errorf("Expected type to be mentioned when scanning, got %q", err)
	}
}

func TestTimestampTZScanUnsupportedFormat(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{`02/03/2001 04:05:06.007 CST`, "ambiguous format"},     // SQL, MDY
		{`03/02/2001 04:05:06.007 CST`, "ambiguous format"},     // SQL, DMY
		{`Sat Feb 03 04:05:06.007 2001 CST`, "not implemented"}, // Postgres, MDY
		{`Sat 03 Feb 04:05:06.007 2001 CST`, "not implemented"}, // Postgres, DMY
		{`03.02.2001 04:05:06.007 CST`, "not implemented"},      // German
	} {
		tstz := TimestampTZ{9, time.Now()}
		err := tstz.Scan(tt.input)

		if err == nil {
			t.Fatal("Expected error for %q, got none", tt.input)
		}

		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}
	}
}

func TestTimestampTZScanTime(t *testing.T) {
	tstz := TimestampTZ{9, time.Now()}
	tt := time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC)
	err := tstz.Scan(tt)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if tstz != (TimestampTZ{Time: tt}) {
		t.Errorf("Expected 2001-02-03 04:05:06.000000007, got %+v", tstz)
	}
}

func BenchmarkTimestampTZScanTime(b *testing.B) {
	var tstz TimestampTZ
	var x, _ interface{} = time.Parse("2006-01-02 15:04:05 MST", `2001-02-03 04:05:06 CET`)
	var y, _ interface{} = time.Parse("2006-01-02 15:04:05 MST", `2001-02-03 04:05:06.007008 CET`)

	for i := 0; i < b.N; i++ {
		tstz.Scan(x)
		tstz.Scan(y)
	}
}

var TimestampTZStringTests = []struct {
	str         string
	timestamptz TimestampTZ
}{
	{`infinity`, TimestampTZ{Infinity: 1}},
	{`-infinity`, TimestampTZ{Infinity: -1}},
	{`2001-02-03 04:05:06-08:09:10`,
		TimestampTZ{Time: time.Date(2001, time.February, 3, 4, 5, 6, 0, time.FixedZone("", -8*3600-9*60-10))}},
	{`2001-02-03 04:05:06.007-08:09:10`,
		TimestampTZ{Time: time.Date(2001, time.February, 3, 4, 5, 6, 7000000, time.FixedZone("", -8*3600-9*60-10))}},
	{`2001-02-03 04:05:06-08:09:10 BC`,
		TimestampTZ{Time: time.Date(-2000, time.February, 3, 4, 5, 6, 0, time.FixedZone("", -8*3600-9*60-10))}},
	{`2001-02-03 04:05:06.007-08:09:10 BC`,
		TimestampTZ{Time: time.Date(-2000, time.February, 3, 4, 5, 6, 7000000, time.FixedZone("", -8*3600-9*60-10))}},
}

func TestTimestampTZScanBytes(t *testing.T) {
	for _, tt := range TimestampTZStringTests {
		bytes := []byte(tt.str)
		tstz := TimestampTZ{9, time.Now()}
		err := tstz.Scan(bytes)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", bytes, err)
		}
		if !reflect.DeepEqual(tstz, tt.timestamptz) {
			t.Errorf("Expected %+v for %q, got %+v", tt.timestamptz, bytes, tstz)
		}
	}
}

func BenchmarkTimestampTZScanBytesISO(b *testing.B) {
	var tstz TimestampTZ
	var w interface{} = []byte(`2001-02-03 04:05:06-09:10:11`)
	var x interface{} = []byte(`2001-02-03 04:05:06.007008-09:10:11`)
	var y interface{} = []byte(`2001-02-03 04:05:06-09:10:11 BC`)
	var z interface{} = []byte(`2001-02-03 04:05:06.007008-09:10:11 BC`)

	for i := 0; i < b.N; i++ {
		tstz.Scan(w)
		tstz.Scan(x)
		tstz.Scan(y)
		tstz.Scan(z)
	}
}

func BenchmarkTimestampTZScanBytesInfinity(b *testing.B) {
	var tstz TimestampTZ
	var x interface{} = []byte(`-infinity`)
	var y interface{} = []byte(`infinity`)

	for i := 0; i < b.N; i++ {
		tstz.Scan(x)
		tstz.Scan(y)
	}
}

func TestTimestampTZScanString(t *testing.T) {
	for _, tt := range TimestampTZStringTests {
		tstz := TimestampTZ{9, time.Now()}
		err := tstz.Scan(tt.str)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.str, err)
		}
		if !reflect.DeepEqual(tstz, tt.timestamptz) {
			t.Errorf("Expected %+v for %q, got %+v", tt.timestamptz, tt.str, tstz)
		}
	}
}

func TestTimestampScanTZError(t *testing.T) {
	now := time.Now()
	tstz := TimestampTZ{9, now}
	err := tstz.Scan("")

	if err == nil {
		t.Error("Expected error, got none")
	}
	if !reflect.DeepEqual(tstz, TimestampTZ{9, now}) {
		t.Errorf("Expected destination not to change, got %+v", tstz)
	}
}
