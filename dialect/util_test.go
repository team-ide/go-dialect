package dialect

import (
	"fmt"
	"reflect"
	"testing"
)

type objString string
type objInt int

func TestString(t *testing.T) {
	fmt.Println("objString:", reflect.TypeOf(objString("xx")).Kind().String())
	fmt.Println("objInt:", reflect.TypeOf(objInt(1)).Kind().String())
}
