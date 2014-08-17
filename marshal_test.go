package parameters

import (
	"strings"
	"testing"
)

type marshalTest struct {
	b []string
	p *TextParameters
}

var marshalTests = []marshalTest{
	{
		b: []string{
			"foo: 3\n",
			"bar: 4\n",
			"baz: 5\n",
		},
		p: &TextParameters{
			"foo": "3",
			"bar": "4",
			"baz": "5",
		},
	},
	{
		// parameter name only
		b: []string{
			"foo\n",
			"bar: 1.4142\n",
		},
		p: &TextParameters{
			"foo": "",
			"bar": "1.4142",
		},
	},
}

func TestMarshal(t *testing.T) {
	for _, test := range marshalTests {
		expects := test.b
		actual := Marshal(test.p)

		for _, expect := range expects {
			if strings.Index(actual, expect) < 0 {
				t.Fatalf(
					"TestMarshal: expect(include) = \"%s\", actual = \"%s\"",
					expect,
					actual,
				)
			}
		}
	}
}
