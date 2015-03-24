package strrand

import (
	"math/rand"
	"strings"
	"time"
)

func makeRange(from rune, to rune) (r []rune) {
	for i := from; i <= to; i++ {
		r = append(r, i)
	}
	return r
}

func concat(ss ...[]rune) (r []rune) {
	for _, s := range ss {
		r = append(r, s...)
	}
	return r
}

type picker []rune

func (p picker) pick() rune {
	// XXX care empty picker
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
	salt = concat(upper, lower, digit, []rune{'.', '/'})
}

var patterns = map[rune]picker{
	'd': picker(digit),
	'D': picker(concat(upper, lower, punct)),
	'w': picker(concat(upper, lower, digit, []rune{'_'})),
	'W': picker(concat(makeRange(33, 47), makeRange(58, 64), makeRange(91, 94), makeRange(96, 96), makeRange(123, 126))),
	's': picker([]rune{' ', '\t'}),
	'S': picker(concat(upper, lower, digit, punct)),

	// these are translated to their double quoted equivalents.
	't': picker([]rune{'\t'}),
	'n': picker([]rune{'\n'}),
	'r': picker([]rune{'\r'}),
	'f': picker([]rune{'\f'}),
	'a': picker([]rune{'\a'}),
	'e': picker([]rune{27}), // escape character
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
		result += string(p.pick())
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
		chr := ([]rune((*chars)[0]))[0]
		*chars = (*chars)[1:]

		switch chr {
		case '\\':
			p, _ := sr.handleEscape(chars)
			pickers = append(pickers, p)
		case '.':
			pickers = append(pickers, any)
		default:
			pickers = append(pickers, picker([]rune{chr}))
		}
	}
	return pickers, nil
}

func (sr *strrand) handleEscape(chars *[]string) (picker, error) {
	chr := ([]rune((*chars)[0]))[0]
	*chars = (*chars)[1:]

	return picker([]rune{chr}), nil
}
