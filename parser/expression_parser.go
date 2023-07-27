package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	p "github.com/a-h/parse"
	m "pasza.org/sr-challenge/model"
)

var primaryParserProxy p.Parser[m.Expr]

func init() {
	primaryParserProxy = primaryParser
}

var lParen = p.Rune('(')
var rParen = p.Rune(')')
var quot = p.Rune('"')

var copyAboveParser = Map(
	p.String("^^"),
	func(string) m.Expr {
		return m.CopyAbove{}
	},
)
var colRefParser = p.RuneInRanges(unicode.Upper)
var intParser = FallibleMap(
	p.OneOrMore[string](p.RuneInRanges(unicode.Digit)),
	func(digits []string) (int, error) {
		return strconv.Atoi(strings.Join(digits, ""))
	},
)

var intLitParser = Map(
	intParser,
	func(value int) m.Expr {
		return m.IntLit(value)
	},
)

// a quoted string, disregard possible escaped quotations for now
var stringLitParser = Map(
	p.SequenceOf3[string, []string, string](quot, p.ZeroOrMore[string](p.RuneNotIn("\"")), quot),
	func(seq p.Tuple3[string, []string, string]) m.Expr {
		return m.StringLit(strings.Join(seq.B, ""))
	},
)

var floatParser = FallibleMap(
	p.SequenceOf3[[]string, string, []string](
		p.OneOrMore[string](p.RuneInRanges(unicode.Digit)),
		p.Rune('.'),
		p.OneOrMore[string](p.RuneInRanges(unicode.Digit)),
	),
	func(seq p.Tuple3[[]string, string, []string]) (float64, error) {
		floatRepr := fmt.Sprintf("%s.%s", strings.Join(seq.A, ""), strings.Join(seq.C, ""))
		return strconv.ParseFloat(floatRepr, 64)
	},
)
var floatLitParser = Map(
	floatParser,
	func(value float64) m.Expr {
		return m.FloatLit(value)
	},
)

var cellRefParser = Map(
	p.SequenceOf2[string, int](colRefParser, intParser),
	func(seq p.Tuple2[string, int]) m.Expr {
		return m.CellRef{
			Col: seq.A,
			Row: seq.B,
		}
	},
)

var copyColumnAboveParser = Map(
	p.SequenceOf2[string, string](colRefParser, p.Rune('^')),
	func(seq p.Tuple2[string, string]) m.Expr {
		return m.CopyColumnAbove{
			Col: seq.A,
		}
	},
)

var copyLastInColumnParser = Map(
	p.SequenceOf3[string, string, string](colRefParser, p.Rune('^'), p.Rune('v')),
	func(seq p.Tuple3[string, string, string]) m.Expr {
		return m.CopyLastInColumn{
			Col: seq.A,
		}
	},
)

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
	p.SequenceOf3[string, int, string](p.Rune('<'), intParser, p.Rune('>')),
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
	func(t p.Tuple3[string, string, int]) m.Expr {
		return m.LabelRelativeRowRef{
			Label:       t.B, // label name
			RelativeRow: t.C, // relative row number
		}
	},
)

var chompWhiteSpace p.Parser[m.NoResult] = Map(
	p.ZeroOrMore(p.RuneIn(" \t")),
	func([]string) m.NoResult {
		return struct{}{}
	},
)

var subExprParser p.Parser[m.Expr] = Map(
	p.SequenceOf3[string, m.Expr, string](
		lParen,
		exprParser,
		rParen,
	),
	func(seq p.Tuple3[string, m.Expr, string]) m.Expr {
		return seq.B
	},
)

// temporary, eventually will include (expr), float, funCall, unaryOps
var primaryParser p.Parser[m.Expr] = p.Any[m.Expr](
	funCallParser,
	stringLitParser,
	floatLitParser,
	intLitParser,
	subExprParser,
	cellRefParser,
	copyAboveParser,
	copyLastInColumnParser,
	copyColumnAboveParser,
	labelRelativeRowRefParser,
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

	pMatch, ok, err := primaryParserProxy.Parse(in)
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
		match, ok, err := primaryParserProxy.Parse(in)
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
	p.SequenceOf3[m.NoResult, string, m.NoResult](
		chompWhiteSpace,
		p.Rune(','),
		chompWhiteSpace,
	),
	func(_ p.Tuple3[m.NoResult, string, m.NoResult]) m.NoResult {
		return m.NoResult{}
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
	func(seq p.Tuple4[string, string, []m.Expr, string]) m.Expr {
		return m.FunCall{
			Name:   seq.A,
			Params: seq.C,
		}
	},
)
