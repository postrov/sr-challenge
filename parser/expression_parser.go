package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	p "github.com/a-h/parse"
	m "pasza.org/sr-challenge/model"
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

var floatLiteral = Map(
	p.SequenceOf3[[]string, string, []string](
		p.OneOrMore[string](p.RuneInRanges(unicode.Digit)),
		p.Rune('.'),
		p.OneOrMore[string](p.RuneInRanges(unicode.Digit)),
	),
	func(seq p.Tuple3[[]string, string, []string]) float64 {
		floatRepr := fmt.Sprintf("%s.%s", strings.Join(seq.A, ""), strings.Join(seq.C, ""))
		res, err := strconv.ParseFloat(floatRepr, 64)
		if err != nil {
			panic("Could not parse float literal")
		}
		return res
	},
)

var cellRef = p.SequenceOf2[string, int](colRef, intLiteral)

// todo: perhaps map to something single
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

// var funcCall = Map(
// 	p.Seq
// )

///// 1 + 2 * 3

var chompWhiteSpace p.Parser[m.NoResult] = Map(
	p.ZeroOrMore(p.RuneIn(" \t")),
	func([]string) m.NoResult {
		return struct{}{}
	},
)

// temporary, eventually will include (expr), float, funCall, unaryOps
var parsePrimary = Map(
	intLiteral,
	func(value int) m.IntLit {
		return m.IntLit(value)
	},
)

var parseBinOp = Map(
	p.RuneIn("+-*/"),
	func(op string) m.BinaryOperator {
		switch op {
		case "+":
			return m.ADD
		case "-":
			return m.SUB
		case "*":
			return m.MUL
		case "/":
			return m.DIV
		default:
			panic("Unknown binary operation")
		}
	},
)

var arExprParser p.Parser[m.Expr] = p.Func(func(in *p.Input) (match m.Expr, ok bool, err error) {
	primaries := make([]m.Expr, 0)
	binOps := make([]m.BinaryOperator, 0)

	chompWhiteSpace.Parse(in)

	pMatch, ok, err := parsePrimary.Parse(in)
	if !ok || err != nil {
		return
	}
	primaries = append(primaries, pMatch)
	for ok {
		chompWhiteSpace.Parse(in)
		oMatch, ok, err := parseBinOp.Parse(in)
		if !ok {
			break
		}
		if err != nil {
			return m.IntLit(3), ok, err
		}
		binOps = append(binOps, oMatch)
		chompWhiteSpace.Parse(in)
		match, ok, err := parsePrimary.Parse(in)
		if !ok || err != nil {
			return match, ok, err
		}
		primaries = append(primaries, match)
	}

	// list of ae1..aeN
	// list of op1..opN-1 (may be empty)
	// match, ok, err = IntLit(3), true, nil
	// match = xxx(primaries, binOps)
	match = fixOperatorPrecedence(primaries, binOps)
	return
})
