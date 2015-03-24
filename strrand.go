package strrand

import (
	"math/rand"
	"strings"
	"time"
)

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

type picker []string

func (p picker) pick() string {
	if len(p) < 1 {
		return ""
	}
	idx := rand.Intn(len(p))
	return p[idx]
}

var upper picker
var lower picker
var digit picker
var punct picker
var any picker
var salt picker
var binary picker

func init() {
	rand.Seed(time.Now().Unix())

	upper = makeRange('A', 'Z')
	lower = makeRange('a', 'z')
	digit = makeRange('0', '9')

	punct = concat(makeRange(33, 47), makeRange(58, 64), makeRange(91, 96), makeRange(123, 126))
	any = concat(upper, lower, digit, punct)
	salt = concat(upper, lower, digit, []string{".", "/"})
}

var patterns = map[string]picker{
	"d": picker(digit),
	"D": picker(concat(upper, lower, punct)),
	"w": picker(concat(upper, lower, digit, []string{"_"})),
	"W": picker(concat(makeRange(33, 47), makeRange(58, 64), makeRange(91, 94), makeRange(96, 96), makeRange(123, 126))),
	"s": picker([]string{" ", "\t"}),
	"S": picker(concat(upper, lower, digit, punct)),

	// these are translated to their double quoted equivalents.
	"t": picker([]string{"\t"}),
	"n": picker([]string{"\n"}),
	"r": picker([]string{"\r"}),
	"f": picker([]string{"\f"}),
	"a": picker([]string{"\a"}),
	"e": picker([]string{string(27)}), // escape character
}

type strrand struct {
}

func New() *strrand {
	return &strrand{}
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
			pickers = append(pickers, picker([]string{chr}))
		}
	}
	return pickers, nil
}

func (sr *strrand) handleEscape(chars *[]string) (picker, error) {
	chr := (*chars)[0]
	*chars = (*chars)[1:]

	return picker([]string{chr}), nil
}
