package parser

import (
	"strconv"
	"strings"
	"unicode"

	p "github.com/a-h/parse"
)

var lParen = p.Rune('(')
var rParen = p.Rune(')')
var quot = p.Rune('"')

var copyAbove = p.String("^^")
var colRef = p.RuneInRanges(unicode.Upper)
var intLiteral = Map(
	p.OneOrMore[string](p.RuneInRanges(unicode.Digit)),
	func(digits []string) int {
		res, err := strconv.Atoi(strings.Join(digits, ""))
		if err != nil {
			panic("Could not parse int literal")
		}

		return res
	},
)

// a quoted string, disregard possible escaped quotations for now
var stringLiteral = Map(
	p.SequenceOf3[string, []string, string](quot, p.ZeroOrMore[string](p.RuneNotIn("\"")), quot),
	func(seq p.Tuple3[string, []string, string]) string {
		return strings.Join(seq.B, "")
	},
)

var cellRef = p.SequenceOf2[string, int](colRef, intLiteral)
var copyLastInColumn = p.SequenceOf2[string, string](colRef, p.Rune('^'))
var labelName = Map(
	p.OneOrMore[string](p.Any[string](
		p.RuneInRanges(unicode.Lower),
		p.RuneIn("_")), // allowed non-letter characters
	),
	func(chars []string) string {
		return strings.Join(chars, "")
	},
)
var relativeRowRef = Map(
	p.SequenceOf3[string, int, string](p.Rune('<'), intLiteral, p.Rune('>')),
	func(t p.Tuple3[string, int, string]) int {
		return t.B
	},
)
var labelRelativeRowRef = Map(
	p.SequenceOf3[string, string, int](
		p.Rune('@'),
		labelName,
		relativeRowRef,
	),
	func(t p.Tuple3[string, string, int]) p.Tuple2[string, int] {
		return p.Tuple2[string, int]{
			A: t.B, // label name
			B: t.C, // relative row number
		}
	},
)

// var funcCall
