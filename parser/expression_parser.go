package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/a-h/parse"
	p "github.com/a-h/parse"
	"pasza.org/sr-challenge/model"
	m "pasza.org/sr-challenge/model"
)

var lParen = p.Rune('(')
var rParen = p.Rune(')')
var quot = p.Rune('"')

var copyAboveParser = p.String("^^")
var colRefParser = p.RuneInRanges(unicode.Upper)
var intLitParser = FallibleMap(
	p.OneOrMore[string](p.RuneInRanges(unicode.Digit)),
	func(digits []string) (int, error) {
		return strconv.Atoi(strings.Join(digits, ""))
	},
)

// a quoted string, disregard possible escaped quotations for now
var stringLiteralParser = Map(
	p.SequenceOf3[string, []string, string](quot, p.ZeroOrMore[string](p.RuneNotIn("\"")), quot),
	func(seq p.Tuple3[string, []string, string]) string {
		return strings.Join(seq.B, "")
	},
)

var floatLitParser = Map(
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

var cellRefParser = p.SequenceOf2[string, int](colRefParser, intLitParser)

// todo: perhaps map to something single
var copyLastInColumnParser = p.SequenceOf2[string, string](colRefParser, p.Rune('^'))
var labelName = Map(
	p.OneOrMore[string](p.Any[string](
		p.RuneInRanges(unicode.Lower),
		p.RuneIn("_")), // allowed non-letter characters
	),
	func(chars []string) string {
		return strings.Join(chars, "")
	},
)
var relativeRowRefParser = Map(
	p.SequenceOf3[string, int, string](p.Rune('<'), intLitParser, p.Rune('>')),
	func(t p.Tuple3[string, int, string]) int {
		return t.B
	},
)
var labelRelativeRowRefParser = Map(
	p.SequenceOf3[string, string, int](
		p.Rune('@'),
		labelName,
		relativeRowRefParser,
	),
	func(t p.Tuple3[string, string, int]) p.Tuple2[string, int] {
		return p.Tuple2[string, int]{
			A: t.B, // label name
			B: t.C, // relative row number
		}
	},
)

var chompWhiteSpace p.Parser[m.NoResult] = Map(
	p.ZeroOrMore(p.RuneIn(" \t")),
	func([]string) m.NoResult {
		return struct{}{}
	},
)

// temporary, eventually will include (expr), float, funCall, unaryOps
var primaryParser = Map(
	intLitParser,
	func(value int) m.IntLit {
		return m.IntLit(value)
	},
)

var binaryOperatorParser = Map(
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

// parses expression composed of arithmetic operations on other expressions
var exprParser p.Parser[m.Expr] = p.Func(func(in *p.Input) (match m.Expr, ok bool, err error) {
	primaries := make([]m.Expr, 0)
	binOps := make([]m.BinaryOperator, 0)

	chompWhiteSpace.Parse(in)

	pMatch, ok, err := primaryParser.Parse(in)
	if !ok || err != nil {
		return
	}
	primaries = append(primaries, pMatch)
	for ok {
		chompWhiteSpace.Parse(in)
		oMatch, ok, err := binaryOperatorParser.Parse(in)
		if !ok {
			break
		}
		if err != nil {
			return m.IntLit(3), ok, err
		}
		binOps = append(binOps, oMatch)
		chompWhiteSpace.Parse(in)
		match, ok, err := primaryParser.Parse(in)
		if !ok || err != nil {
			return match, ok, err
		}
		primaries = append(primaries, match)
	}

	// list of ae1..aeN
	// list of op1..opN-1 (may be empty)
	match = fixOperatorPrecedence(primaries, binOps)
	return
})

var funNameParser = Map(
	p.OneOrMore(p.RuneInRanges(unicode.Letter)),
	func(chars []string) string {
		return strings.Join(chars, "")
	},
)

var argSeparatorParser = Map(
	parse.SequenceOf3[model.NoResult, string, model.NoResult](
		chompWhiteSpace,
		parse.Rune(','),
		chompWhiteSpace,
	),
	func(_ parse.Tuple3[model.NoResult, string, model.NoResult]) model.NoResult {
		return model.NoResult{}
	},
)

var argListParser = SeparatedList0[m.Expr, m.NoResult](exprParser, argSeparatorParser)

var funCallParser = Map(
	p.SequenceOf4[string, string, []m.Expr, string](
		funNameParser,
		lParen,
		argListParser,
		rParen,
	),
	func(seq p.Tuple4[string, string, []m.Expr, string]) m.FunCall {
		return m.FunCall{
			Name:   seq.A,
			Params: seq.C,
		}
	},
)
