package strrand

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMakeRange(t *testing.T) {
	s := makeRange('A', 'C')

	if !reflect.DeepEqual(s, []string{"A", "B", "C"}) {
		t.Errorf("something wrong")
	}

	ss, _ := New().Randregex(`..\.adf.!\d\w`)
	fmt.Println(ss)
}

func TestPicker(t *testing.T) {
	p := chrPicker([]string{"X"})

	if !reflect.DeepEqual(p, chrPicker([]string{"X"})) {
		t.Errorf("something wrong")
	}

	if p.pick() != "X" {
		t.Errorf("something wrong")
	}
}

func TestVariantPicker(t *testing.T) {
	vp := variantPicker{
		min:    1,
		max:    3,
		picker: chrPicker([]string{"a", "b", "c", "d"}),
	}

	r := vp.pick()
	fmt.Println(r)

	if len(r) < 1 || len(r) > 3 {
		t.Errorf("something wrong")
	}
}
