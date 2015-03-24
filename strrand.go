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

var upper []rune
var lower []rune
var digit []rune
var punct []rune
var any []rune
var salt []rune
var binary []rune

func init() {
	rand.Seed(time.Now().Unix())

	upper = makeRange('A', 'Z')
	lower = makeRange('a', 'z')
	digit = makeRange('0', '9')

	punct = concat(makeRange(33, 47), makeRange(58, 64), makeRange(91, 96), makeRange(123, 126))
	any = concat(upper, lower, digit, punct)
	salt = concat(upper, lower, digit, []rune{'.', '/'})
}

var patterns = map[rune]([]rune){
	'd': digit,
	'D': concat(upper, lower, punct),
	'w': concat(upper, lower, digit, []rune{'_'}),
	'W': concat(makeRange(33, 47), makeRange(58, 64), makeRange(91, 94), makeRange(96, 96), makeRange(123, 126)),
	's': []rune{' ', '\t'},
	'S': concat(upper, lower, digit, punct),

	// these are translated to their double quoted equivalents.
	't': []rune{'\t'},
	'n': []rune{'\n'},
	'r': []rune{'\r'},
	'f': []rune{'\f'},
	'a': []rune{'\a'},
	'e': []rune{27}, // escape character
}

type strrand struct {
}

func New() *strrand {
	return &strrand{}
}

func (sr *strrand) Randregex(pattern string) (string, error) {
	chars := func() *[]string {
		c := strings.Split(pattern, "")
		return &c
	}()

	result := ""
	for len(*chars) > 0 {
		chr := (*chars)[0][0]
		*chars = (*chars)[1:]

		switch chr {
		case '\\':
			result += sr.handleEscape(chars)
		case '.':
			result += sr.handleDot()
		default:
			result += string(chr)
		}
	}

	return result, nil
}

func randGet(s []rune) string {
	idx := rand.Intn(len(s))
	return string(s[idx])
}

func (sr *strrand) handleDot() string {
	return randGet(any)
}

func (st *strrand) handleEscape(chars *[]string) string {
	chr := (*chars)[0]
	*chars = (*chars)[1:]
	return chr
}
