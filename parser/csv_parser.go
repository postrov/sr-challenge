package parser

import (
	"log"
	"strings"

	p "github.com/a-h/parse"
	"pasza.org/sr-challenge/model"
)

// var newLine = p.SequenceOf2(p.Optional(p.Rune('\r')), p.Rune('\n'))
var colDelimiter = p.Rune('|')
var rowDelimiter = p.NewLine
var rawCell p.Parser[string] = p.Func(func(in *p.Input) (match string, ok bool, err error) {
	match, ok, err = p.StringUntil(p.Any(colDelimiter, rowDelimiter, p.EOF[string]())).Parse(in)
	p.Any(colDelimiter).Parse(in)
	return
})

var cell = Map(rawCell, func(s string) model.Cell {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "=") {
		// parse formula
		return model.FormulaCell{
			Formula: model.DummyFormula{
				RawValue: "formula: " + s[1:],
			},
		}
	} else {
		return model.StringCell{
			Value: s,
		}
		// just string value
	}
})

var row p.Parser[[]model.Cell] = p.Func(func(in *p.Input) (match []model.Cell, ok bool, err error) {
	match, ok, err = p.UntilEOF(cell, rowDelimiter).Parse(in)
	rowDelimiter.Parse(in)
	return
})

var CSV p.Parser[[][]model.Cell] = p.Until(row, p.EOF[string]())

const _csvData string = `!date         |!transaction_id                        |!tokens        |!token_prices          |!total_cost
2022-02-20    |=concat("t_", text(incFrom(1)))        |btc,eth,dai    |38341.88,2643.77,1.0003|=sum(spread(split(D2, ",")))
2022-02-21    |=^^                                    |bch,eth,dai    |304.38,2621.15,1.0001  |=E^+sum(spread(split(D3, ",")))
2022-02-22    |=^^                                    |sol,eth,dai    |85,2604.17,0.9997      |=^^



!fee          |!cost_threshold                        |               |                       |
0.09          |10000                                  |               |                       |



!adjusted_cost|                                       |               |                       |
=E^v+(E^v*A9) |                                       |               |                       |

!cost_too_high|                                       |               |                       |
=text(bte(@adjusted_cost<1>, @cost_threshold<1>)      |               |                       |`

const csvData string = `ala | =ma | kota
kot | ma | =ale`

func Parser() string {
	input := p.NewInput(csvData)
	log.Print("got input")

	match, ok, err := CSV.Parse(input)

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
			log.Printf("  %s\n", c)
		}
	}

	// this is for row
	// for i, c := range match {
	// 	log.Printf("cell[%d]: %s", i, c)
	// }

	return "xxx"
}
