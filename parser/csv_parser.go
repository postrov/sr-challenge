package parser

import (
	"log"
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

var stringCellParser = Map(
	p.StringUntilEOF[m.NoResult](
		p.Func[m.NoResult](func(in *p.Input) (m.NoResult, bool, error) {
			return m.NoResult{}, false, nil
		}),
	),
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

var CSVParser p.Parser[[][]m.Cell] = p.Until(rowParser, p.EOF[string]())

const csvData string = `!date         |!transaction_id                        |!tokens        |!token_prices          |!total_cost
2022-02-20    |=concat("t_", text(incFrom(1)))        |btc,eth,dai    |38341.88,2643.77,1.0003|=sum(spread(split(D2, ",")))
2022-02-21    |=^^                                    |bch,eth,dai    |304.38,2621.15,1.0001  |=E^+sum(spread(split(D3, ",")))
2022-02-22    |=^^                                    |sol,eth,dai    |85,2604.17,0.9997      |=^^



!fee          |!cost_threshold                        |               |                       |
0.09          |10000                                  |               |                       |



!adjusted_cost|                                       |               |                       |
=E^v+(E^v*A9) |                                       |               |                       |

!cost_too_high|                                       |               |                       |
=text(bte(@adjusted_cost<1>, @cost_threshold<1>)      |               |                       |`

func Parser() string {
	input := p.NewInput(csvData)

	match, ok, err := CSVParser.Parse(input)

	if err != nil {
		log.Fatalf("failed to parse: %v\n", err)
	}

	if !ok {
		log.Print("expected CSV not matched\n")
	}

	// this is for whole CSV
	for i, r := range match {
		log.Printf("row: %d\n", i)
		for _, c := range r {
			log.Print(c)
		}
	}

	// this is for row
	// for i, c := range match {
	// 	log.Printf("cell[%d]: %s", i, c)
	// }

	return "xxx"
}
