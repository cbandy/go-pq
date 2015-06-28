package pq

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

func TestBoolArrayScanUnsupported(t *testing.T) {
	var arr BoolArray
	err := arr.Scan(1)

	if err == nil {
		t.Fatal("Expected error when scanning from int")
	}
	if !strings.Contains(err.Error(), "int to BoolArray") {
		t.Errorf("Expected type to be mentioned when scanning, got %q", err)
	}
}

var BoolArrayStringTests = []struct {
	str string
	arr BoolArray
}{
	{`{}`, BoolArray{}},
	{`{t}`, BoolArray{true}},
	{`{f,t}`, BoolArray{false, true}},
}

func TestBoolArrayScanBytes(t *testing.T) {
	for _, tt := range BoolArrayStringTests {
		bytes := []byte(tt.str)
		arr := BoolArray{true, true, true}
		err := arr.Scan(bytes)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", bytes, err)
		}
		if !reflect.DeepEqual(arr, tt.arr) {
			t.Errorf("Expected %+v for %q, got %+v", tt.arr, bytes, arr)
		}
	}
}

func BenchmarkBoolArrayScanAllocated(b *testing.B) {
	var a BoolArray = make(BoolArray, 10)
	var x interface{} = []byte(`{t,f,t,f,t,f,t,f,t,f}`)

	for i := 0; i < b.N; i++ {
		a.Scan(x)
	}
}

func BenchmarkBoolArrayScanBytes(b *testing.B) {
	var a BoolArray
	var x interface{} = []byte(`{t,f,t,f,t,f,t,f,t,f}`)

	for i := 0; i < b.N; i++ {
		a = BoolArray{}
		a.Scan(x)
	}
}

func TestBoolArrayScanString(t *testing.T) {
	for _, tt := range BoolArrayStringTests {
		arr := BoolArray{true, true, true}
		err := arr.Scan(tt.str)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.str, err)
		}
		if !reflect.DeepEqual(arr, tt.arr) {
			t.Errorf("Expected %+v for %q, got %+v", tt.arr, tt.str, arr)
		}
	}
}

func TestBoolArrayScanError(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{``, "unexpected format"},
		{`{`, "unexpected format"},
		{`}`, "unexpected format"},
		{`{{}`, "unexpected format"},
		{`{}}`, "unexpected format"},
		{`{x}`, "unexpected format"},
		{`{,}`, "unexpected format"},
		{`{t,}`, "unexpected format"},
		{`{,t}`, "unexpected format"},
	} {
		arr := BoolArray{true, true, true}
		err := arr.Scan(tt.input)

		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}
		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}
		if !strings.Contains(err.Error(), tt.input) {
			t.Errorf("Expected error to contain %q, got %q", tt.input, err)
		}
		if !reflect.DeepEqual(arr, BoolArray{true, true, true}) {
			t.Errorf("Expected destination not to change for %q, got %+v", tt.input, arr)
		}
	}
}

func TestBoolArrayValue(t *testing.T) {
	result, err := BoolArray(nil).Value()

	if err != nil {
		t.Fatalf("Expected no error for nil, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil, got %q", result)
	}

	result, err = BoolArray([]bool{}).Value()

	if err != nil {
		t.Fatalf("Expected no error for empty, got %v", err)
	}
	if expected := `{}`; !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected empty, got %q", result)
	}

	result, err = BoolArray([]bool{false, true, false}).Value()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if expected := `{f,t,f}`; !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func BenchmarkBoolArrayValue(b *testing.B) {
	rand.Seed(1)
	x := make([]bool, 10)
	for i := 0; i < len(x); i++ {
		x[i] = rand.Intn(2) == 0
	}
	a := BoolArray(x)

	for i := 0; i < b.N; i++ {
		a.Value()
	}
}

func TestByteaArrayScanUnsupported(t *testing.T) {
	var arr ByteaArray
	err := arr.Scan(1)

	if err == nil {
		t.Fatal("Expected error when scanning from int")
	}
	if !strings.Contains(err.Error(), "int to ByteaArray") {
		t.Errorf("Expected type to be mentioned when scanning, got %q", err)
	}
}

var ByteaArrayStringTests = []struct {
	str string
	arr ByteaArray
}{
	{`{}`, ByteaArray{}},
	{`{"\\xfeff"}`, ByteaArray{{'\xFE', '\xFF'}}},
	{`{"\\xdead","\\xbeef"}`, ByteaArray{{'\xDE', '\xAD'}, {'\xBE', '\xEF'}}},
}

func TestByteaArrayScanBytes(t *testing.T) {
	for _, tt := range ByteaArrayStringTests {
		bytes := []byte(tt.str)
		arr := ByteaArray{{2}, {6}, {0, 0}}
		err := arr.Scan(bytes)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", bytes, err)
		}
		if !reflect.DeepEqual(arr, tt.arr) {
			t.Errorf("Expected %+v for %q, got %+v", tt.arr, bytes, arr)
		}
	}
}

func BenchmarkByteaArrayScanAllocated(b *testing.B) {
	var a ByteaArray = make(ByteaArray, 10)
	var x interface{} = []byte(`{"\\xfe","\\xff","\\xdead","\\xbeef","\\xfe","\\xff","\\xdead","\\xbeef","\\xfe","\\xff"}`)

	for i := 0; i < b.N; i++ {
		a.Scan(x)
	}
}

func BenchmarkByteaArrayScanBytes(b *testing.B) {
	var a ByteaArray
	var x interface{} = []byte(`{"\\xfe","\\xff","\\xdead","\\xbeef","\\xfe","\\xff","\\xdead","\\xbeef","\\xfe","\\xff"}`)

	for i := 0; i < b.N; i++ {
		a = ByteaArray{}
		a.Scan(x)
	}
}

func TestByteaArrayScanString(t *testing.T) {
	for _, tt := range ByteaArrayStringTests {
		arr := ByteaArray{{2}, {6}, {0, 0}}
		err := arr.Scan(tt.str)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.str, err)
		}
		if !reflect.DeepEqual(arr, tt.arr) {
			t.Errorf("Expected %+v for %q, got %+v", tt.arr, tt.str, arr)
		}
	}
}

func TestByteaArrayScanError(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{``, "unexpected format"},
		{`{`, "unexpected format"},
		{`}`, "unexpected format"},
		{`{{}`, "unexpected format"},
		{`{}}`, "unexpected format"},
		{`{x}`, "unexpected format"},
		{`{,}`, "unexpected format"},
		{`{"\\xfeff",}`, "unexpected format"},
		{`{,"\\xfeff"}`, "unexpected format"},
	} {
		arr := ByteaArray{{2}, {6}, {0, 0}}
		err := arr.Scan(tt.input)

		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}
		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}
		if !strings.Contains(err.Error(), tt.input) {
			t.Errorf("Expected error to contain %q, got %q", tt.input, err)
		}
		if !reflect.DeepEqual(arr, ByteaArray{{2}, {6}, {0, 0}}) {
			t.Errorf("Expected destination not to change for %q, got %+v", tt.input, arr)
		}
	}
}

func TestByteaArrayValue(t *testing.T) {
	result, err := ByteaArray(nil).Value()

	if err != nil {
		t.Fatalf("Expected no error for nil, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil, got %q", result)
	}

	result, err = ByteaArray([][]byte{}).Value()

	if err != nil {
		t.Fatalf("Expected no error for empty, got %v", err)
	}
	if expected := `{}`; !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected empty, got %q", result)
	}

	result, err = ByteaArray([][]byte{{'\xDE', '\xAD', '\xBE', '\xEF'}, {'\xFE', '\xFF'}, {}}).Value()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if expected := `{"\\xdeadbeef","\\xfeff","\\x"}`; !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func BenchmarkByteaArrayValue(b *testing.B) {
	rand.Seed(1)
	x := make([][]byte, 10)
	for i := 0; i < len(x); i++ {
		x[i] = make([]byte, len(x))
		for j := 0; j < len(x); j++ {
			x[i][j] = byte(rand.Int())
		}
	}
	a := ByteaArray(x)

	for i := 0; i < b.N; i++ {
		a.Value()
	}
}

func TestFloat64ArrayScanUnsupported(t *testing.T) {
	var arr Float64Array
	err := arr.Scan(true)

	if err == nil {
		t.Fatal("Expected error when scanning from bool")
	}
	if !strings.Contains(err.Error(), "bool to Float64Array") {
		t.Errorf("Expected type to be mentioned when scanning, got %q", err)
	}
}

var Float64ArrayStringTests = []struct {
	str string
	arr Float64Array
}{
	{`{}`, Float64Array{}},
	{`{1.2}`, Float64Array{1.2}},
	{`{3.456,7.89}`, Float64Array{3.456, 7.89}},
	{`{3,1,2}`, Float64Array{3, 1, 2}},
}

func TestFloat64ArrayScanBytes(t *testing.T) {
	for _, tt := range Float64ArrayStringTests {
		bytes := []byte(tt.str)
		arr := Float64Array{5, 5, 5}
		err := arr.Scan(bytes)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", bytes, err)
		}
		if !reflect.DeepEqual(arr, tt.arr) {
			t.Errorf("Expected %+v for %q, got %+v", tt.arr, bytes, arr)
		}
	}
}

func BenchmarkFloat64ArrayScanAllocated(b *testing.B) {
	var a Float64Array = make(Float64Array, 10)
	var x interface{} = []byte(`{1.2,3.4,5.6,7.8,9.01,2.34,5.67,8.90,1.234,5.678}`)

	for i := 0; i < b.N; i++ {
		a.Scan(x)
	}
}

func BenchmarkFloat64ArrayScanBytes(b *testing.B) {
	var a Float64Array
	var x interface{} = []byte(`{1.2,3.4,5.6,7.8,9.01,2.34,5.67,8.90,1.234,5.678}`)

	for i := 0; i < b.N; i++ {
		a = Float64Array{}
		a.Scan(x)
	}
}

func TestFloat64ArrayScanString(t *testing.T) {
	for _, tt := range Float64ArrayStringTests {
		arr := Float64Array{5, 5, 5}
		err := arr.Scan(tt.str)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.str, err)
		}
		if !reflect.DeepEqual(arr, tt.arr) {
			t.Errorf("Expected %+v for %q, got %+v", tt.arr, tt.str, arr)
		}
	}
}

func TestFloat64ArrayScanError(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{``, "unexpected format"},
		{`{`, "unexpected format"},
		{`}`, "unexpected format"},
		{`{{}`, "unexpected format"},
		{`{}}`, "unexpected format"},
		{`{x}`, "unexpected format"},
		{`{,}`, "unexpected format"},
		{`{1.2,}`, "unexpected format"},
		{`{,1.2}`, "unexpected format"},
	} {
		arr := Float64Array{5, 5, 5}
		err := arr.Scan(tt.input)

		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}
		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}
		if !strings.Contains(err.Error(), tt.input) {
			t.Errorf("Expected error to contain %q, got %q", tt.input, err)
		}
		if !reflect.DeepEqual(arr, Float64Array{5, 5, 5}) {
			t.Errorf("Expected destination not to change for %q, got %+v", tt.input, arr)
		}
	}
}

func TestFloat64ArrayValue(t *testing.T) {
	result, err := Float64Array(nil).Value()

	if err != nil {
		t.Fatalf("Expected no error for nil, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil, got %q", result)
	}

	result, err = Float64Array([]float64{}).Value()

	if err != nil {
		t.Fatalf("Expected no error for empty, got %v", err)
	}
	if expected := `{}`; !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected empty, got %q", result)
	}

	result, err = Float64Array([]float64{1.2, 3.4, 5.6}).Value()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if expected := `{1.2,3.4,5.6}`; !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func BenchmarkFloat64ArrayValue(b *testing.B) {
	rand.Seed(1)
	x := make([]float64, 10)
	for i := 0; i < len(x); i++ {
		x[i] = rand.NormFloat64()
	}
	a := Float64Array(x)

	for i := 0; i < b.N; i++ {
		a.Value()
	}
}

func TestInt64ArrayScanUnsupported(t *testing.T) {
	var arr Int64Array
	err := arr.Scan(true)

	if err == nil {
		t.Fatal("Expected error when scanning from bool")
	}
	if !strings.Contains(err.Error(), "bool to Int64Array") {
		t.Errorf("Expected type to be mentioned when scanning, got %q", err)
	}
}

var Int64ArrayStringTests = []struct {
	str string
	arr Int64Array
}{
	{`{}`, Int64Array{}},
	{`{12}`, Int64Array{12}},
	{`{345,678}`, Int64Array{345, 678}},
}

func TestInt64ArrayScanBytes(t *testing.T) {
	for _, tt := range Int64ArrayStringTests {
		bytes := []byte(tt.str)
		arr := Int64Array{5, 5, 5}
		err := arr.Scan(bytes)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", bytes, err)
		}
		if !reflect.DeepEqual(arr, tt.arr) {
			t.Errorf("Expected %+v for %q, got %+v", tt.arr, bytes, arr)
		}
	}
}

func BenchmarkInt64ArrayScanAllocated(b *testing.B) {
	var a Int64Array = make(Int64Array, 10)
	var x interface{} = []byte(`{1,2,3,4,5,6,7,8,9,0}`)

	for i := 0; i < b.N; i++ {
		a.Scan(x)
	}
}

func BenchmarkInt64ArrayScanBytes(b *testing.B) {
	var a Int64Array
	var x interface{} = []byte(`{1,2,3,4,5,6,7,8,9,0}`)

	for i := 0; i < b.N; i++ {
		a = Int64Array{}
		a.Scan(x)
	}
}

func TestInt64ArrayScanString(t *testing.T) {
	for _, tt := range Int64ArrayStringTests {
		arr := Int64Array{5, 5, 5}
		err := arr.Scan(tt.str)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.str, err)
		}
		if !reflect.DeepEqual(arr, tt.arr) {
			t.Errorf("Expected %+v for %q, got %+v", tt.arr, tt.str, arr)
		}
	}
}

func TestInt64ArrayScanError(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{``, "unexpected format"},
		{`{`, "unexpected format"},
		{`}`, "unexpected format"},
		{`{{}`, "unexpected format"},
		{`{}}`, "unexpected format"},
		{`{x}`, "unexpected format"},
		{`{,}`, "unexpected format"},
		{`{1,}`, "unexpected format"},
		{`{,1}`, "unexpected format"},
	} {
		arr := Int64Array{5, 5, 5}
		err := arr.Scan(tt.input)

		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}
		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}
		if !strings.Contains(err.Error(), tt.input) {
			t.Errorf("Expected error to contain %q, got %q", tt.input, err)
		}
		if !reflect.DeepEqual(arr, Int64Array{5, 5, 5}) {
			t.Errorf("Expected destination not to change for %q, got %+v", tt.input, arr)
		}
	}
}

func TestInt64ArrayValue(t *testing.T) {
	result, err := Int64Array(nil).Value()

	if err != nil {
		t.Fatalf("Expected no error for nil, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil, got %q", result)
	}

	result, err = Int64Array([]int64{}).Value()

	if err != nil {
		t.Fatalf("Expected no error for empty, got %v", err)
	}
	if expected := `{}`; !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected empty, got %q", result)
	}

	result, err = Int64Array([]int64{1, 2, 3}).Value()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if expected := `{1,2,3}`; !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func BenchmarkInt64ArrayValue(b *testing.B) {
	rand.Seed(1)
	x := make([]int64, 10)
	for i := 0; i < len(x); i++ {
		x[i] = rand.Int63()
	}
	a := Int64Array(x)

	for i := 0; i < b.N; i++ {
		a.Value()
	}
}

func TestStringArrayScanUnsupported(t *testing.T) {
	var arr StringArray
	err := arr.Scan(true)

	if err == nil {
		t.Fatal("Expected error when scanning from bool")
	}
	if !strings.Contains(err.Error(), "bool to StringArray") {
		t.Errorf("Expected type to be mentioned when scanning, got %q", err)
	}
}

var StringArrayStringTests = []struct {
	str string
	arr StringArray
}{
	{`{}`, StringArray{}},
	{`{t}`, StringArray{"t"}},
	{`{f,1}`, StringArray{"f", "1"}},
	{`{"a\\b","c d",","}`, StringArray{"a\\b", "c d", ","}},
}

func TestStringArrayScanBytes(t *testing.T) {
	for _, tt := range StringArrayStringTests {
		bytes := []byte(tt.str)
		arr := StringArray{"x", "x", "x"}
		err := arr.Scan(bytes)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", bytes, err)
		}
		if !reflect.DeepEqual(arr, tt.arr) {
			t.Errorf("Expected %+v for %q, got %+v", tt.arr, bytes, arr)
		}
	}
}

func BenchmarkStringArrayScanAllocated(b *testing.B) {
	var a StringArray = make(StringArray, 10)
	var x interface{} = []byte(`{a,b,c,d,e,f,g,h,i,j}`)

	for i := 0; i < b.N; i++ {
		a.Scan(x)
	}
}

func BenchmarkStringArrayScanBytes(b *testing.B) {
	var a StringArray
	var x interface{} = []byte(`{a,b,c,d,e,f,g,h,i,j}`)
	var y interface{} = []byte(`{"\a","\b","\c","\d","\e","\f","\g","\h","\i","\j"}`)

	for i := 0; i < b.N; i++ {
		a = StringArray{}
		a.Scan(x)
		a = StringArray{}
		a.Scan(y)
	}
}

func TestStringArrayScanString(t *testing.T) {
	for _, tt := range StringArrayStringTests {
		arr := StringArray{"x", "x", "x"}
		err := arr.Scan(tt.str)

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.str, err)
		}
		if !reflect.DeepEqual(arr, tt.arr) {
			t.Errorf("Expected %+v for %q, got %+v", tt.arr, tt.str, arr)
		}
	}
}

func TestStringArrayScanError(t *testing.T) {
	for _, tt := range []struct {
		input, err string
	}{
		{``, "unexpected format"},
		{`{`, "unexpected format"},
		{`}`, "unexpected format"},
		{`{{}`, "unexpected format"},
		{`{}}`, "unexpected format"},
		{`{"}`, "unexpected format"},
		{`{,}`, "unexpected format"},
		{`{a,}`, "unexpected format"},
		{`{,a}`, "unexpected format"},
		{`{"\}`, "unexpected format"},
		{`{"\"}`, "unexpected format"},
	} {
		arr := StringArray{"x", "x", "x"}
		err := arr.Scan(tt.input)

		if err == nil {
			t.Fatalf("Expected error for %q, got none", tt.input)
		}
		if !strings.Contains(err.Error(), tt.err) {
			t.Errorf("Expected error to contain %q for %q, got %q", tt.err, tt.input, err)
		}
		if !strings.Contains(err.Error(), tt.input) {
			t.Errorf("Expected error to contain %q, got %q", tt.input, err)
		}
		if !reflect.DeepEqual(arr, StringArray{"x", "x", "x"}) {
			t.Errorf("Expected destination not to change for %q, got %+v", tt.input, arr)
		}
	}
}

func TestStringArrayValue(t *testing.T) {
	result, err := StringArray(nil).Value()

	if err != nil {
		t.Fatalf("Expected no error for nil, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil, got %q", result)
	}

	result, err = StringArray([]string{}).Value()

	if err != nil {
		t.Fatalf("Expected no error for empty, got %v", err)
	}
	if expected := `{}`; !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected empty, got %q", result)
	}

	result, err = StringArray([]string{`a`, `\b`, `c"`, `d,e`}).Value()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if expected := `{"a","\\b","c\"","d,e"}`; !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func BenchmarkStringArrayValue(b *testing.B) {
	x := make([]string, 10)
	for i := 0; i < len(x); i++ {
		x[i] = strings.Repeat(`abc"def\ghi`, 5)
	}
	a := StringArray(x)

	for i := 0; i < b.N; i++ {
		a.Value()
	}
}

func TestGenericArrayValueUnsupported(t *testing.T) {
	_, err := GenericArray{true}.Value()

	if err == nil {
		t.Fatal("Expected error for bool")
	}
	if !strings.Contains(err.Error(), "bool to array") {
		t.Errorf("Expected type to be mentioned, got %q", err)
	}
}

type ByteArrayValuer [1]byte
type ByteSliceValuer []byte
type FuncArrayValuer struct {
	delimiter func() string
	value     func() (driver.Value, error)
}

func (a ByteArrayValuer) Value() (driver.Value, error) { return a[:], nil }
func (b ByteSliceValuer) Value() (driver.Value, error) { return []byte(b), nil }
func (f FuncArrayValuer) ArrayDelimiter() string       { return f.delimiter() }
func (f FuncArrayValuer) Value() (driver.Value, error) { return f.value() }

func TestGenericArrayValue(t *testing.T) {
	result, err := GenericArray{nil}.Value()

	if err != nil {
		t.Fatalf("Expected no error for nil, got %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil, got %q", result)
	}

	Tilde := func(v driver.Value) FuncArrayValuer {
		return FuncArrayValuer{
			func() string { return "~" },
			func() (driver.Value, error) { return v, nil }}
	}

	for _, tt := range []struct {
		result string
		input  interface{}
	}{
		{`{}`, []bool{}},
		{`{true}`, []bool{true}},
		{`{true,false}`, []bool{true, false}},
		{`{true,false}`, [2]bool{true, false}},

		{`{}`, [][]int{{}}},
		{`{}`, [][]int{{}, {}}},
		{`{{1}}`, [][]int{{1}}},
		{`{{1},{2}}`, [][]int{{1}, {2}}},
		{`{{1,2},{3,4}}`, [][]int{{1, 2}, {3, 4}}},
		{`{{1,2},{3,4}}`, [2][2]int{{1, 2}, {3, 4}}},

		{`{"a","\\b","c\"","d,e"}`, []string{`a`, `\b`, `c"`, `d,e`}},
		{`{"a","\\b","c\"","d,e"}`, [][]byte{{'a'}, {'\\', 'b'}, {'c', '"'}, {'d', ',', 'e'}}},

		{`{NULL}`, []*int{nil}},
		{`{0,NULL}`, []*int{new(int), nil}},

		{`{NULL}`, []sql.NullString{{}}},
		{`{"\"",NULL}`, []sql.NullString{{`"`, true}, {}}},

		{`{"a","b"}`, []ByteArrayValuer{{'a'}, {'b'}}},
		{`{{"a","b"},{"c","d"}}`, [][]ByteArrayValuer{{{'a'}, {'b'}}, {{'c'}, {'d'}}}},

		{`{"e","f"}`, []ByteSliceValuer{{'e'}, {'f'}}},
		{`{{"e","f"},{"g","h"}}`, [][]ByteSliceValuer{{{'e'}, {'f'}}, {{'g'}, {'h'}}}},

		{`{1~2}`, []FuncArrayValuer{Tilde(int64(1)), Tilde(int64(2))}},
		{`{{1~2}~{3~4}}`, [][]FuncArrayValuer{{Tilde(int64(1)), Tilde(int64(2))}, {Tilde(int64(3)), Tilde(int64(4))}}},
	} {
		result, err := GenericArray{tt.input}.Value()

		if err != nil {
			t.Fatalf("Expected no error for %q, got %v", tt.input, err)
		}
		if !reflect.DeepEqual(result, tt.result) {
			t.Errorf("Expected %q for %q, got %q", tt.result, tt.input, result)
		}
	}
}

func TestGenericArrayValueErrors(t *testing.T) {
	var v []interface{}

	v = []interface{}{func() {}}
	if _, err := (GenericArray{v}).Value(); err == nil {
		t.Errorf("Expected error for %q, got nil", v)
	}

	v = []interface{}{nil, func() {}}
	if _, err := (GenericArray{v}).Value(); err == nil {
		t.Errorf("Expected error for %q, got nil", v)
	}
}

func BenchmarkGenericArrayValueBools(b *testing.B) {
	rand.Seed(1)
	x := make([]bool, 10)
	for i := 0; i < len(x); i++ {
		x[i] = rand.Intn(2) == 0
	}
	a := GenericArray{x}

	for i := 0; i < b.N; i++ {
		a.Value()
	}
}

func BenchmarkGenericArrayValueFloat64s(b *testing.B) {
	rand.Seed(1)
	x := make([]float64, 10)
	for i := 0; i < len(x); i++ {
		x[i] = rand.NormFloat64()
	}
	a := GenericArray{x}

	for i := 0; i < b.N; i++ {
		a.Value()
	}
}

func BenchmarkGenericArrayValueInt64s(b *testing.B) {
	rand.Seed(1)
	x := make([]int64, 10)
	for i := 0; i < len(x); i++ {
		x[i] = rand.Int63()
	}
	a := GenericArray{x}

	for i := 0; i < b.N; i++ {
		a.Value()
	}
}

func BenchmarkGenericArrayValueByteSlices(b *testing.B) {
	x := make([][]byte, 10)
	for i := 0; i < len(x); i++ {
		x[i] = bytes.Repeat([]byte(`abc"def\ghi`), 5)
	}
	a := GenericArray{x}

	for i := 0; i < b.N; i++ {
		a.Value()
	}
}

func BenchmarkGenericArrayValueStrings(b *testing.B) {
	x := make([]string, 10)
	for i := 0; i < len(x); i++ {
		x[i] = strings.Repeat(`abc"def\ghi`, 5)
	}
	a := GenericArray{x}

	for i := 0; i < b.N; i++ {
		a.Value()
	}
}

func TestArrayValueBackend(t *testing.T) {
	db := openTestConn(t)
	defer db.Close()

	for _, tt := range []struct {
		s string
		v driver.Valuer
	}{
		{`ARRAY[true, false]`, BoolArray([]bool{true, false})},
		{`ARRAY[E'\\xdead', E'\\xbeef']`, ByteaArray([][]byte{{'\xDE', '\xAD'}, {'\xBE', '\xEF'}})},
		{`ARRAY[1.2, 3.4]`, Float64Array([]float64{1.2, 3.4})},
		{`ARRAY[1, 2, 3]`, Int64Array([]int64{1, 2, 3})},
		{`ARRAY['a', E'\\b', 'c"', 'd,e']`, StringArray([]string{`a`, `\b`, `c"`, `d,e`})},
	} {
		var x int
		err := db.QueryRow(`SELECT 1 WHERE `+tt.s+` <> $1`, tt.v).Scan(&x)
		if err != sql.ErrNoRows {
			t.Errorf("Expected %v to equal %q, got %q", tt.v, tt.s, err)
		}
	}
}
