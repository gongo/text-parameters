package parameters

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	var actual bytes.Buffer
	encoder := NewEncoder(&actual)

	type st struct {
		TypeFloat32 float32
		TypeFloat64 float64
		TypeInt     int
		TypeString  string
		TypeUint    uint
		typeInt     int // unexported field
	}
	s := st{
		TypeFloat32: 1.4142,
		TypeFloat64: 2.8888,
		TypeInt:     -20,
		TypeString:  "GoLang",
		TypeUint:    20,
		typeInt:     -2000,
	}

	encoder.Encode(&s)
	expect := `TypeFloat32: 1.4142
TypeFloat64: 2.8888
TypeInt: -20
TypeString: GoLang
TypeUint: 20
`

	if actual.String() != expect {
		t.Fatalf("TestEncode: expect = %s, actual = %s", expect, actual.String())
	}
}

func TestEncodeWithTag(t *testing.T) {
	var actual bytes.Buffer
	encoder := NewEncoder(&actual)

	type st struct {
		Received   int     `textparam:"packet-received"`
		Time       float64 `textparam:"transfer-time"`
		Extra      string
		extraField string `textparam:"extra-transfer-time"`
	}
	s := st{
		Received:   1024,
		Time:       12.345,
		Extra:      "extra-field",
		extraField: "unexported extra-field",
	}

	encoder.Encode(&s)
	expect := `Extra: extra-field
packet-received: 1024
transfer-time: 12.345
`

	if actual.String() != expect {
		t.Fatalf("TestEncode: expect = %s, actual = %s", expect, actual.String())
	}
}
