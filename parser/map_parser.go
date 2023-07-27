package parser

import p "github.com/a-h/parse"

type mapParser[A, B any] struct {
	A      p.Parser[A]
	mapper func(A) B
}

func (p mapParser[A, B]) Parse(in *p.Input) (match B, ok bool, err error) {
	matchA, ok, err := p.A.Parse(in)
	if err != nil || !ok {
		return
	}

	match = p.mapper(matchA)
	return
}

// Map takes a parser, a mapping function and returns new parser that applies the function to result of first parser
func Map[A, B any](a p.Parser[A], mapper func(A) B) p.Parser[B] {
	return mapParser[A, B]{
		A:      a,
		mapper: mapper,
	}
}

// / fallible map
type fallibleMapParser[A, B any] struct {
	A      p.Parser[A]
	mapper func(A) (B, error)
}

func (p fallibleMapParser[A, B]) Parse(in *p.Input) (match B, ok bool, err error) {
	matchA, ok, err := p.A.Parse(in)
	if err != nil || !ok {
		return
	}

	match, err = p.mapper(matchA)
	return
}

func FallibleMap[A, B any](a p.Parser[A], mapper func(A) (b B, err error)) p.Parser[B] {
	return fallibleMapParser[A, B]{
		A:      a,
		mapper: mapper,
	}
}
