package parameters

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	received := 10
	time := 0.3838
	var money uint = 1980
	body := fmt.Sprintf("Received: %d\nTime: %f\nMoney: %d", received, time, money)

	type st struct {
		Received int
		Time     float64
		Money    uint
	}
	s := st{}

	decoder := NewDecorder(strings.NewReader(body))

	if err := decoder.Decode(&s); err != nil {
		t.Fatal(err)
	}

	if s.Received != received {
		t.Fatalf("TestDecode: expect = %d, actual = %v", received, s.Received)
	}

	if s.Time != time {
		t.Fatalf("TestDecode: expect = %f, actual = %v", time, s.Time)
	}

	if s.Money != money {
		t.Fatalf("TestDecode: expect = %f, actual = %v", money, s.Money)
	}
}

func TestDecodeWithTag(t *testing.T) {
	received := 10
	time := 0.3838
	body := fmt.Sprintf(
		"packet-received: %d\ntransfer-time: %f\nextra: foobar\n",
		received,
		time,
	)

	type st struct {
		Received   int     `parameters:"packet-received"`
		Time       float64 `parameters:"transfer-time"`
		Extra      string
		extraField string `parameters:"transfer-time"`
	}
	s := st{}

	decoder := NewDecorder(strings.NewReader(body))

	if err := decoder.Decode(&s); err != nil {
		t.Fatal(err)
	}

	if s.Received != received {
		t.Fatalf("TestDecode: expect = %d, actual = `%v`", received, s.Received)
	}

	if s.Time != time {
		t.Fatalf("TestDecode: expect = %f, actual = `%v`", time, s.Time)
	}

	// not store to `Extra` field because has not tag.
	if s.Extra != "" {
		t.Fatalf("TestDecode: expect = \"\", actual = `%v`", s.Extra)
	}

	// not store to `extraField` field because that is unexported field.
	if s.extraField != "" {
		t.Fatalf("TestDecode: expect = \"\", actual = `%v`", s.extraField)
	}
}

type decodeFailedTest struct {
	b string
	s interface{}
	e error
}

var decodeFailedTests = []decodeFailedTest{
	{
		b: "A: foobar\n",
		s: &struct {
			A int // should be string
		}{},
		e: &DecodeFieldTypeError{},
	},
	{
		b: "A: foobar\nB: 1.4142\n",
		s: &struct {
			A string
			B int // should be float{32,64}
		}{},
		e: &DecodeFieldTypeError{},
	},
	{
		b: "C: -1\n",
		s: &struct {
			C uint32 // should be signed integer
		}{},
		e: &DecodeFieldTypeError{},
	},
	{
		b: "ddd: 12.345\neeeee: hoge",
		s: &struct {
			D3 float64 `parameters:"ddd"`
			E5 float64 `parameters:"eeeee"` // should be string
		}{},
		e: &DecodeFieldTypeError{},
	},
	{
		b: "d dd: 12.345\n", // format error
		s: &struct{}{},
		e: &DecodeFormatError{},
	},
	{
		b: "ddd: 12.345\n",
		s: struct{}{},
		e: &CodingStructPointerError{},
	},
}

func TestDecodeFailed(t *testing.T) {
	for _, test := range decodeFailedTests {
		decoder := NewDecorder(strings.NewReader(test.b))
		err := decoder.Decode(test.s)

		if err == nil {
			t.Fatalf(
				"TestDecodedFailed: Error should occur (b = \"%s\", s = %v)",
				test.b,
				test.s,
			)
		}

		actualError := reflect.TypeOf(err)
		expectError := reflect.TypeOf(test.e)

		if actualError != expectError {
			t.Fatalf(
				"TestDecodedFailed: expect: %s, actual = %s (%s)",
				expectError,
				actualError,
				err,
			)
		}
	}
}
