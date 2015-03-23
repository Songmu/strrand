package strrand

import (
	"math/rand"
	"strings"
	"time"
)

func makeRange(from rune, to rune) []rune {
	r := []rune{}
	for i := from; i <= to; i++ {
		r = append(r, i)
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

	punct = makeRange(33, 47)
	punctB := makeRange(58, 64)
	punctC := makeRange(123, 126)
	punct = append(punct, punctB...)
	punct = append(punct, punctC...)

	any = append(upper, lower...)
	any = append(any, digit...)
	any = append(any, punct...)

	salt = append(upper, lower...)
	salt = append(salt, digit...)
	salt = append(salt, '.', '/')

	binary = makeRange(0, 255)
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
