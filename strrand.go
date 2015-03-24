package strrand

import (
	"math/rand"
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

var upper chrPicker
var lower chrPicker
var digit chrPicker
var punct chrPicker
var any chrPicker
var salt chrPicker
var binary chrPicker

func init() {
	rand.Seed(time.Now().Unix())

	upper = makeRange('A', 'Z')
	lower = makeRange('a', 'z')
	digit = makeRange('0', '9')

	punct = concat(makeRange(33, 47), makeRange(58, 64), makeRange(91, 96), makeRange(123, 126))
	any = concat(upper, lower, digit, punct)
	salt = concat(upper, lower, digit, []string{".", "/"})
}

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
	pickers, _ := sr.doRandregex(pattern)

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
			p, _ := sr.handleEscape(chars)
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
	chr := (*chars)[0]
	*chars = (*chars)[1:]

	return chrPicker([]string{chr}), nil
}
