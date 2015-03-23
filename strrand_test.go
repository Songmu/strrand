package strrand

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMakeRange(t *testing.T) {
	s := makeRange('A', 'C')

	if !reflect.DeepEqual(s, []rune{'A', 'B', 'C'}) {
		t.Errorf("something wrong")
	}

	ss, _ := New().Randregex("..\\.adf.\\!")
	fmt.Println(ss)
}
