package parser

import (
	"strings"
	"testing"
	"unicode"

	"github.com/a-h/parse"
	"github.com/stretchr/testify/assert"
	"pasza.org/sr-challenge/model"
)

var nameParser = Map(
	parse.OneOrMore(parse.RuneInRanges(unicode.Letter)),
	func(chars []string) string {
		return strings.Join(chars, "")
	},
)

var separatorParser = Map(
	parse.SequenceOf3[model.NoResult, string, model.NoResult](
		chompWhiteSpace,
		parse.Rune(','),
		chompWhiteSpace,
	),
	func(_ parse.Tuple3[model.NoResult, string, model.NoResult]) model.NoResult {
		return model.NoResult{}
	},
)

func TestSeparatedListParser1(t *testing.T) {
	separatedListParser := SeparatedList1[string, model.NoResult](nameParser, separatorParser)

	cases := []struct {
		in   string
		want []string
	}{
		{"ala,ma,kota", []string{"ala", "ma", "kota"}},
		{"ala, ma , kota", []string{"ala", "ma", "kota"}},
	}

	for _, c := range cases {
		input := parse.NewInput(c.in)
		match, ok, err := separatedListParser.Parse(input)
		assert.Nil(t, err)
		assert.True(t, ok)
		assert.Equal(t, c.want, match)
	}
}

func TestSeparatedListParser0(t *testing.T) {
	separatedListParser := SeparatedList0[string, model.NoResult](nameParser, separatorParser)

	cases := []struct {
		in   string
		want []string
	}{
		{"kot,ma , ale", []string{"kot", "ma", "ale"}},
		{"", nil},
	}

	for _, c := range cases {
		input := parse.NewInput(c.in)
		match, ok, err := separatedListParser.Parse(input)
		assert.Nil(t, err)
		assert.True(t, ok)
		assert.Equal(t, c.want, match)
	}
}
