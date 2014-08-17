package parameters_test

import (
	"fmt"
	"strings"

	"github.com/gongo/text-parameters"
)

func ExampleMarshal() {
	params := &parameters.TextParameters{
		"Foo": "foobar",
		"Bar": "golang",
	}

	encoded := parameters.Marshal(params)
	fmt.Println(encoded)
	// Output:
	// Bar: golang
	// Foo: foobar
}

func ExampleUnmarshal() {
	body := `foo: 3
bar: 4
baz: 5`

	obj, _ := parameters.Unmarshal(strings.NewReader(body))

	fmt.Println(obj.Get("foo"), obj.Get("bar"), obj.Get("baz"))
	// Output:
	// 3 4 5
}

func ExampleDecoder() {
	body := `Foo: 3.14
barbaz: golang
Piyo: 123`

	s := struct {
		Foo  float64
		Bar  string `textparam:"barbaz"`
		Piyo int
	}{}

	decoder := parameters.NewDecorder(strings.NewReader(body))
	decoder.Decode(&s)

	fmt.Println("s.Foo  =", s.Foo)
	fmt.Println("s.Bar  =", s.Bar)
	fmt.Println("s.Piyo =", s.Piyo)
	// Output:
	// s.Foo  = 3.14
	// s.Bar  = golang
	// s.Piyo = 123
}
