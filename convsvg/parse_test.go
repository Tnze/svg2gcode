package convsvg

import (
	"testing"
)

func TestFieldsD(t *testing.T) {
	fs, err := fieldsD(`M5016 1595 c-123 -31 -266 -129 -266 -182 0 -56 102 -112 250 -138
90 -16 369 -30 399 -20 19 6 21 14 21 74 0 125 -45 211 -133 258 -44 23 -193
27 -271 8z`)
	if err != nil {
		t.Error(err)
	}
	t.Log(fs)
	//for _, v := range fs {
	//	switch v.(type) {
	//	case byte, rune:
	//		t.Logf("%T(%c)", v, v)
	//	default:
	//		t.Logf("%v", v)
	//	}
	//}
}
