package parser

import (
	"github.com/a-h/parse"
	p "github.com/a-h/parse"
)

type separatedListParser[A, B any] struct {
	elementParser   p.Parser[A]
	separatorParser p.Parser[B]
	emptyOk         bool
}

func (p separatedListParser[A, B]) Parse(in *parse.Input) (match []A, ok bool, err error) {
	matchA, ok, err := p.elementParser.Parse(in)
	if err != nil {
		return nil, false, err
	}

	if !ok {
		return nil, p.emptyOk, nil
	}

	// got first match
	result := make([]A, 0, 1)
	result = append(result, matchA)

	// try matching more elements
	for ok {
		var match parse.Tuple2[B, A]
		match, ok, err = parse.SequenceOf2[B, A](p.separatorParser, p.elementParser).Parse(in)
		if err != nil {
			return nil, false, err
		}
		if ok {
			result = append(result, match.B)
		}
	}

	return result, true, nil
}

// Matches one or more `a` parses, separated by `b` parses
func SeparatedList1[A, B any](a parse.Parser[A], b parse.Parser[B]) parse.Parser[[]A] {
	return separatedListParser[A, B]{
		elementParser:   a,
		separatorParser: b,
		emptyOk:         false,
	}
}

// Matches zero or more `a` parses, separated by `b` parses
func SeparatedList0[A, B any](a parse.Parser[A], b parse.Parser[B]) parse.Parser[[]A] {
	return separatedListParser[A, B]{
		elementParser:   a,
		separatorParser: b,
		emptyOk:         true,
	}
}
