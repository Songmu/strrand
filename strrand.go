package strrand

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type picker interface {
	pick() string
}

func makeRange(from rune, to rune) (r []string) {
	for i := from; i <= to; i++ {
		r = append(r, string(i))
	}
	return r
}

func concat(ss ...[]string) (r []string) {
	for _, s := range ss {
		r = append(r, s...)
	}
	return r
}

type chrPicker []string

func (p chrPicker) pick() string {
	if len(p) < 1 {
		return ""
	}
	idx := rand.Intn(len(p))
	return p[idx]
}

type variantPicker struct {
	min    uint
	max    uint
	picker picker
}

func (vp variantPicker) pick() (str string) {
	num := vp.decidePickNum()
	for i := 0; i < num; i++ {
		str += vp.picker.pick()
	}
	return str
}

func (vp variantPicker) decidePickNum() int {
	if vp.max < vp.min {
		return 0
	}
	fluct := vp.max - vp.min
	return int(vp.min) + rand.Intn(int(fluct))
}

func init() {
	rand.Seed(time.Now().Unix())
}

var upper chrPicker = makeRange('A', 'Z')
var lower chrPicker = makeRange('a', 'z')
var digit chrPicker = makeRange('0', '9')

var punct chrPicker = concat(makeRange(33, 47), makeRange(58, 64), makeRange(91, 96), makeRange(123, 126))
var any chrPicker = concat(upper, lower, digit, punct)
var salt chrPicker = concat(upper, lower, digit, []string{".", "/"})

var patterns = map[string]chrPicker{
	"d": chrPicker(digit),
	"D": chrPicker(concat(upper, lower, punct)),
	"w": chrPicker(concat(upper, lower, digit, []string{"_"})),
	"W": chrPicker(concat(makeRange(33, 47), makeRange(58, 64), makeRange(91, 94), makeRange(96, 96), makeRange(123, 126))),
	"s": chrPicker([]string{" ", "\t"}),
	"S": chrPicker(concat(upper, lower, digit, punct)),

	// these are translated to their double quoted equivalents.
	"t": chrPicker([]string{"\t"}),
	"n": chrPicker([]string{"\n"}),
	"r": chrPicker([]string{"\r"}),
	"f": chrPicker([]string{"\f"}),
	"a": chrPicker([]string{"\a"}),
	"e": chrPicker([]string{string(27)}), // escape character
}

type strrand struct {
	Max uint
}

func New() *strrand {
	return &strrand{}
}

var defaultMax uint = 10

func (sr *strrand) max() uint {
	if sr.Max < 1 {
		return defaultMax
	}
	return sr.Max
}

func (sr *strrand) Randregex(pattern string) (string, error) {
	result := ""
	pickers, err := sr.doRandregex(pattern)
	if err != nil {
		return result, err
	}

	for _, p := range pickers {
		result += p.pick()
	}
	return result, nil
}

func (sr *strrand) doRandregex(pattern string) ([]picker, error) {
	pickers := []picker{}
	chars := func() *[]string {
		c := strings.Split(pattern, "")
		return &c
	}()

	for len(*chars) > 0 {
		chr := (*chars)[0]
		*chars = (*chars)[1:]

		switch chr {
		case "\\":
			p, err := sr.handleEscape(chars)
			if err != nil {
				return []picker{}, err
			}
			pickers = append(pickers, p)
		case ".":
			pickers = append(pickers, any)
		default:
			pickers = append(pickers, chrPicker([]string{chr}))
		}
	}
	return pickers, nil
}

func (sr *strrand) handleEscape(chars *[]string) (chrPicker, error) {
	if len(*chars) < 1 {
		return chrPicker([]string{}), fmt.Errorf("regex not terminated")
	}
	chr := (*chars)[0]
	*chars = (*chars)[1:]

	if val, ok := patterns[chr]; ok {
		return val, nil
	}

	switch chr {
	case "x":
		if len(*chars) < 2 {
			return chrPicker([]string{}), fmt.Errorf("invalid hex format")
		}
		hex := (*chars)[0] + (*chars)[1]
		*chars = (*chars)[2:]

		c, err := strconv.ParseInt(hex, 16, 32)
		if err != nil {
			return chrPicker([]string{}), err
		}
		return chrPicker([]string{string(c)}), nil
	default:
		return chrPicker([]string{chr}), nil
	}
}
