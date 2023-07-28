package parser

import (
	"strings"

	p "github.com/a-h/parse"
	m "pasza.org/sr-challenge/model"
)

// var newLine = p.SequenceOf2(p.Optional(p.Rune('\r')), p.Rune('\n'))
var colDelimiterParser = p.Rune('|')
var rowDelimiterParser = p.NewLine
var rawCellParser p.Parser[string] = p.Func(func(in *p.Input) (match string, ok bool, err error) {
	match, ok, err = p.StringUntil(p.Any(colDelimiterParser, rowDelimiterParser, p.EOF[string]())).Parse(in)
	p.Any(colDelimiterParser).Parse(in)
	return
})

var intCellParser = Map(
	p.SequenceOf2[int, m.NoResult](
		intParser,
		p.EOF[m.NoResult](),
	),
	func(seq p.Tuple2[int, m.NoResult]) m.Cell {
		return m.IntCell{
			Value: seq.A,
		}
	},
)
var floatCellParser = Map(
	p.SequenceOf2[float64, m.NoResult](
		floatParser,
		p.EOF[m.NoResult](),
	),
	func(seq p.Tuple2[float64, m.NoResult]) m.Cell {
		return m.FloatCell{
			Value: seq.A,
		}
	},
)

var labelCellParser = Map(
	p.SequenceOf2[string, string](
		p.Rune('!'),
		stringUntilEOFParser,
	),
	func(seq p.Tuple2[string, string]) m.Cell {
		return m.LabelCell{
			Label: seq.B,
		}
	},
)

var formulaCellParser = Map(
	p.SequenceOf3[string, m.Expr](
		p.Rune('='),
		exprParser,
		p.EOF[m.NoResult](),
	),
	func(seq p.Tuple3[string, m.Expr, m.NoResult]) m.Cell {
		return m.FormulaCell{
			Formula: seq.B,
		}
	},
)

var stringUntilEOFParser = p.StringUntilEOF[m.NoResult](
	p.Func[m.NoResult](func(in *p.Input) (m.NoResult, bool, error) {
		return m.NoResult{}, false, nil
	}),
)

var stringCellParser = Map(
	stringUntilEOFParser,
	func(s string) m.Cell {
		return m.StringCell{
			Value: s,
		}
	},
)

var cellParser = FallibleMap(
	rawCellParser,
	func(s string) (m.Cell, error) {
		s = strings.TrimSpace(s)
		match, ok, err := p.Any[m.Cell](
			labelCellParser,
			formulaCellParser,
			floatCellParser,
			intCellParser,
			stringCellParser,
		).Parse(p.NewInput(s))
		if err != nil {
			return nil, err
		}
		if !ok {
			panic("Failed to parse cell")
		}
		return match, nil
	},
)

var rowParser p.Parser[[]m.Cell] = p.Func(func(in *p.Input) (match []m.Cell, ok bool, err error) {
	match, ok, err = p.UntilEOF(cellParser, rowDelimiterParser).Parse(in)
	rowDelimiterParser.Parse(in)
	return
})

var csvParser p.Parser[[][]m.Cell] = p.Until(rowParser, p.EOF[string]())

func ParseCSV(csvData string) ([][]m.Cell, bool, error) {
	input := p.NewInput(csvData)

	return csvParser.Parse(input)
}
