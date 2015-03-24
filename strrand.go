package strrand

import (
	"fmt"
	"math/rand"
	"regexp"
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
	min    uint64
	max    uint64
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
	if vp.max == vp.min {
		return int(vp.min)
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
	// W: punct without '_'
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
		case "[":
			p, err := sr.handleBracket(chars)
			if err != nil {
				return []picker{}, err
			}
			pickers = append(pickers, p)
		case "*":
			*chars = append(strings.Split("{0,}", ""), (*chars)...)
		case "+":
			*chars = append(strings.Split("{1,}", ""), (*chars)...)
		case "?":
			*chars = append(strings.Split("{0,1}", ""), (*chars)...)
		case "{":
			if len(pickers) < 1 {
				return pickers, fmt.Errorf("invalid bracket")
			}
			pop := pickers[len(pickers)-1]
			pickers = pickers[:len(pickers)-1]
			p, err := sr.handleBrace(chars, pop)
			if err != nil {
				return pickers, err
			}
			pickers = append(pickers, p)
		default:
			pickers = append(pickers, chrPicker([]string{chr}))
		}
	}
	return pickers, nil
}

func (sr *strrand) handleEscape(chars *[]string) (picker, error) {
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
			return nil, fmt.Errorf("invalid hex format")
		}
		hex := (*chars)[0] + (*chars)[1]
		*chars = (*chars)[2:]

		c, err := strconv.ParseInt(hex, 16, 32)
		if err != nil {
			return nil, err
		}
		return chrPicker([]string{string(c)}), nil
	default:
		return chrPicker([]string{chr}), nil
	}
}

func (sr *strrand) handleBracket(chars *[]string) (picker, error) {
	p := chrPicker([]string{})
	escaped := false
	rangeJustFinished := false
	startChr := ""

	for {
		if len(*chars) < 1 {
			return nil, fmt.Errorf("unmatched []")
		}
		chr := (*chars)[0]
		*chars = (*chars)[1:]

		if !escaped {
			if chr == "]" {
				break
			}

			if chr == `\` {
				escaped = true
				continue
			}

			if chr == "-" && len(p) > 0 {
				if rangeJustFinished {
					return nil, fmt.Errorf("invalid range")
				}
				startChr = p[len(p)-1]
				p = p[:len(p)-1]
				continue
			}
		}

		escaped = false
		if startChr != "" {
			r := makeRange(([]rune(startChr))[0], ([]rune(chr))[0])
			p = append(p, r...)
			rangeJustFinished = true
			startChr = ""
			continue
		}
		rangeJustFinished = false
		p = append(p, chr)
	}
	return p, nil
}

var rangeChrRe = regexp.MustCompile(`[0-9,]`)
var rangeRe = regexp.MustCompile(`^([0-9]+),([0-9]+)?$`)

func (sr *strrand) handleBrace(chars *[]string, p picker) (picker, error) {
	hasComma := false
	buf := ""
	for {
		if len(*chars) < 1 {
			return nil, fmt.Errorf("bad range %s", buf)
		}
		chr := (*chars)[0]
		*chars = (*chars)[1:]
		if chr == "}" {
			break
		}
		if !rangeChrRe.MatchString(chr) {
			return nil, fmt.Errorf("'%s' inside {} not supported", chr)
		}
		if chr == "," {
			hasComma = true
		}
		buf += chr
	}
	var min, max uint64

	if hasComma {
		matches := rangeRe.FindStringSubmatch(buf)
		if len(matches) == 0 {
			return nil, fmt.Errorf("malformed range: {%s}", buf)
		}

		min, _ = strconv.ParseUint(matches[1], 10, 64)
		if matches[2] == "" {
			max = uint64(sr.max())
		} else {
			max, _ = strconv.ParseUint(matches[2], 10, 64)
		}
		if max < min {
			return nil, fmt.Errorf("bad range {%s}", buf)
		}
	} else {
		min, _ = strconv.ParseUint(buf, 10, 64)
		max = min
	}

	return variantPicker{
		max:    max,
		min:    min,
		picker: p,
	}, nil
}
