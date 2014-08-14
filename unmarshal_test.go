package parameters

import (
	"reflect"
	"strings"
	"testing"
)

type unmarshalTest struct {
	b string
	p *TextParameters
}

var unmarshalTests = []unmarshalTest{
	{
		b: `foo: 3
bar: 4
baz: 5
`,
		p: &TextParameters{
			"foo": "3",
			"bar": "4",
			"baz": "5",
		},
	},
	{
		// Spaces that are not unified
		b: `foo   :    gold
bar:           silver
baz              :bronze
`,
		p: &TextParameters{
			"foo": "gold",
			"bar": "silver",
			"baz": "bronze",
		},
	},
	{
		// parameter name only
		b: `foo
bar: 1.4142
baz
piyo: GoLang
`,
		p: &TextParameters{
			"foo":  "",
			"bar":  "1.4142",
			"baz":  "",
			"piyo": "GoLang",
		},
	},
	{
		b: `F:oo: 3`,
		p: &TextParameters{
			"F": "oo: 3",
		},
	},
	{
		b: `(^^): 3`, // %x21-39 or %x3B-7E
		p: &TextParameters{
			"(^^)": "3",
		},
	},
}

func TestUnmarshal(t *testing.T) {
	for _, test := range unmarshalTests {
		expect := test.p
		actual, err := Unmarshal(strings.NewReader(test.b))

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(expect, actual) {
			t.Fatalf(
				"TestUnmarshal: expect = %v, actual = %v",
				*expect,
				*actual,
			)
		}
	}
}

var unmarshalFailedTests = []unmarshalTest{
	{
		b: `f oo: 3`, // can't use space in parameter name
		p: &TextParameters{},
	},
	{
		b: `愛: 3`, // can't use non-ascii
		p: &TextParameters{},
	},
}

func TestUnmarshalFailed(t *testing.T) {
	for _, test := range unmarshalFailedTests {
		_, err := Unmarshal(strings.NewReader(test.b))

		switch err.(type) {
		case *DecodeFormatError:
			// expect
		default:
			if err == nil {
				t.Fatalf(
					"TestUnmarshalFailed: Error should occur (b: `%s`, p: %v)",
					test.b,
					*test.p,
				)
			}

			t.Fatal(err)
		}
	}
}
