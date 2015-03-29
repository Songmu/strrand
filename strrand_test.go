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
}

func TestCreateGenerator(t *testing.T) {
	reg := `.\.a{1,4}[a-c][A-B]?\d*\s+\![-OMG!\?]`
	g, _ := New().CreateGenerator(reg)

	if !reflect.DeepEqual(g, pickers([]picker{
		any,
		chrPicker([]string{"."}),
		variantPicker{
			min:    1,
			max:    4,
			picker: chrPicker([]string{"a"}),
		},
		chrPicker([]string{"a", "b", "c"}),
		variantPicker{
			min:    0,
			max:    1,
			picker: chrPicker([]string{"A", "B"}),
		},
		variantPicker{
			min:    0,
			max:    10,
			picker: digit,
		},
		variantPicker{
			min:    1,
			max:    10,
			picker: chrPicker([]string{" ", "\t"}),
		},
		chrPicker([]string{"!"}),
		chrPicker([]string{"-", "O", "M", "G", "!", "?"}),
	})) {
		t.Errorf("error")
	}

	fmt.Printf("`%s`: %s\n", reg, g.Generate())
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

func TestRandomString(t *testing.T) {
	str, _ := RandomString("hoge")

	if str != "hoge" {
		t.Errorf("something wrong")
	}
}

func TestCreateGeneratorNaruho(t *testing.T) {
	reg := `なるほど?`
	g, _ := New().CreateGenerator(reg)

	if !reflect.DeepEqual(g, pickers([]picker{
		chrPicker([]string{"な"}),
		chrPicker([]string{"る"}),
		chrPicker([]string{"ほ"}),
		variantPicker{
			min:    0,
			max:    1,
			picker: chrPicker([]string{"ど"}),
		},
	})) {
		t.Errorf("error")
	}

	fmt.Printf("`%s`: %s\n", reg, g.Generate())
}
