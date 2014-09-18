package parameters

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	type st struct {
		Received int
		Time     float64
		Money    uint
	}
	actual := st{}
	expect := st{
		Received: 10,
		Time:     0.3838,
		Money:    uint(1980),
	}

	body := fmt.Sprintf(
		"Received: %d\nTime: %f\nMoney: %d",
		expect.Received,
		expect.Time,
		expect.Money,
	)
	decoder := NewDecorder(strings.NewReader(body))

	if err := decoder.Decode(&actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("TestDecode: expect = %v, actual = %v", expect, actual)
	}
}

func TestDecodeWithTag(t *testing.T) {
	type st struct {
		Received   int     `parameters:"packet-received"`
		Time       float64 `parameters:"transfer-time"`
		Extra      string
		extraField string `parameters:"transfer-time"`
	}
	actual := st{}
	expect := st{
		Received:   10,
		Time:       0.3838,
		Extra:      "", // Has not tag and same name parameter in body.
		extraField: "", // Unexported field.
	}

	body := fmt.Sprintf(
		"packet-received: %d\ntransfer-time: %f\nextra: foobar\n",
		expect.Received,
		expect.Time,
	)
	decoder := NewDecorder(strings.NewReader(body))

	if err := decoder.Decode(&actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expect) {
		t.Fatalf("TestDecodeWithTag: expect = %v, actual = %v", expect, actual)
	}
}

type decodeFailedTest struct {
	b string
	e error
	s interface{}
}

var decodeFailedTests = []decodeFailedTest{
	{
		b: "A: foobar\n",
		e: &DecodeFieldTypeError{},
		s: &struct {
			A int // should be string
		}{},
	},
	{
		b: "A: foobar\nB: 1.4142\n",
		e: &DecodeFieldTypeError{},
		s: &struct {
			A string
			B int // should be float{32,64}
		}{},
	},
	{
		b: "C: -1\n",
		e: &DecodeFieldTypeError{},
		s: &struct {
			C uint32 // should be signed integer
		}{},
	},
	{
		b: "ddd: 12.345\neeeee: hoge",
		e: &DecodeFieldTypeError{},
		s: &struct {
			D3 float64 `parameters:"ddd"`
			E5 float64 `parameters:"eeeee"` // should be string
		}{},
	},
	{
		b: "d dd: 12.345\n",
		e: &DecodeFormatError{},
		s: &struct{}{},
	},
	{
		b: "ddd: 12.345\n",
		e: &CodingStructPointerError{},
		s: struct{}{}, // should be a pointer
	},
}

func TestDecodeFailedFieldType(t *testing.T) {
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

		expectErrorType := reflect.TypeOf(test.e)
		actualErrorType := reflect.TypeOf(err)

		if expectErrorType != actualErrorType {
			t.Fatalf(
				"TestDecodedFailed: expect: %s, actual = %s (%s)",
				expectErrorType,
				actualErrorType,
				err,
			)
		}
	}
}
