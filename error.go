package parameters

import (
	"fmt"
	"reflect"
)

type DecodeFormatError struct {
	invalidLine string
}

func (e *DecodeFormatError) Error() string {
	return fmt.Sprintf("text-parameters: DecodeFormatError(\"%s\")", e.invalidLine)
}

type DecodeFieldTypeError struct {
	t    reflect.StructField
	body string
}

func (e *DecodeFieldTypeError) Error() string {
	return fmt.Sprintf("text-parameters: DecodeFieldTypeError(\"%s\" into `%s %s`)",
		e.body, e.t.Name, e.t.Type)
}
