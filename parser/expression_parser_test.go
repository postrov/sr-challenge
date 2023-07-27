package parser

import (
	"fmt"
	"math"
	"testing"

	p "github.com/a-h/parse"
	"github.com/stretchr/testify/assert"
	m "pasza.org/sr-challenge/model"
)

func TestLabelNameParser(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"label1234", "label"},
		{"token_price", "token_price"},
	}

	for _, c := range cases {
		input := p.NewInput(c.in)
		match, ok, err := labelName.Parse(input)
		assert.True(t, ok)
		assert.Nil(t, err)
		assert.Equal(t, c.want, match)
	}
}

func TestRelativeRowRefParser(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"<5>", 5},
		{"<1234>", 1234},
	}

	for _, c := range cases {
		input := p.NewInput(c.in)
		match, ok, err := relativeRowRefParser.Parse(input)
		assert.True(t, ok)
		assert.Nil(t, err)
		assert.Equal(t, c.want, match)
	}
}

const epsilon float64 = 0.000001

func TestFloatLiteralParser(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"31.337", 31.337},
	}

	for _, c := range cases {
		input := p.NewInput(c.in)
		match, ok, err := floatLitParser.Parse(input)

		assert.True(t, ok)
		assert.Nil(t, err)

		value, ok := match.(m.FloatLit)
		assert.True(t, ok)
		assert.Less(t, math.Abs(c.want-float64(value)), epsilon)
	}
}

func TestLabelRelativeRowRefParser(t *testing.T) {
	cases := []struct {
		in            string
		wantLabelName string
		wantRowNum    int
	}{
		{"@token_price<77>", "token_price", 77},
	}

	for _, c := range cases {
		input := p.NewInput(c.in)
		match, ok, err := labelRelativeRowRefParser.Parse(input)
		assert.True(t, ok)
		assert.Nil(t, err)

		value, ok := match.(m.LabelRelativeRowRef)
		assert.True(t, ok)

		assert.Equal(t, c.wantLabelName, value.Label)
		assert.Equal(t, c.wantRowNum, value.RelativeRow)
	}
}

func TestStringLiteralParserOk(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{`"a quoted string"`, "a quoted string"},
		{`"a quoted string"and then some`, "a quoted string"},
		{`""`, ""},
	}

	for _, c := range cases {
		input := p.NewInput(c.in)
		match, ok, err := stringLitParser.Parse(input)
		assert.True(t, ok)
		assert.Nil(t, err)
		value, ok := match.(m.StringLit)
		assert.True(t, ok)
		assert.Equal(t, c.want, string(value))
	}
}

func TestStringLiteralParserFail(t *testing.T) {
	cases := []struct {
		in string
		ok bool
	}{
		{`15`, false},
	}

	for _, c := range cases {
		input := p.NewInput(c.in)
		_, ok, err := stringLitParser.Parse(input)
		assert.Equal(t, c.ok, ok)
		assert.Nil(t, err)
	}
}

func TestPrimaryParser(t *testing.T) {
	input := p.NewInput("15 + 33")
	match, ok, err := primaryParser.Parse(input)

	assert.True(t, ok)
	assert.Nil(t, err)
	fmt.Print(match)
}

/// arExpr

func TestExprParser(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"15", "15"},
		{"1 + 2 * 3", "1 + (2 * 3)"},
		{"1 + 2 * 3 + 4", "(1 + (2 * 3)) + 4"},
		{"1 + 2 * 3 + 4 * 5 - 6", "((1 + (2 * 3)) + (4 * 5)) - 6"},
		{`E^+sum(spread(split(D3, ",")))`, `E^ + sum(spread(split(D3, ",")))`},
	}
	for _, c := range cases {
		input := p.NewInput(c.in)
		match, ok, err := exprParser.Parse(input)
		assert.True(t, ok)
		assert.Nil(t, err)
		matchRepr := fmt.Sprint(match)
		assert.Equal(t, c.want, matchRepr)
	}
}

func TestIntLitParser(t *testing.T) {
	in := p.NewInput("123")
	match, ok, err := primaryParser.Parse(in)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, m.IntLit(123), match)

	in = p.NewInput("1 + 2 * 3")
	match, ok, err = primaryParser.Parse(in)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, m.IntLit(1), match)
}

func TestFunCallParser(t *testing.T) {
	cases := []struct {
		in           string
		wantName     string
		wantArgc     int
		wantArgvRepr []string
	}{
		{
			"add(1, 2 + 3, 4)",
			"add",
			3,
			[]string{"1", "2 + 3", "4"},
		},
	}

	for _, c := range cases {
		input := p.NewInput(c.in)
		match, ok, err := funCallParser.Parse(input)
		assert.True(t, ok)
		assert.Nil(t, err)

		value, ok := match.(m.FunCall)
		assert.True(t, ok)

		assert.Equal(t, c.wantName, value.Name)
		assert.Equal(t, c.wantArgc, len(value.Params))
		for i, v := range value.Params {
			assert.Equalf(t, c.wantArgvRepr[i], fmt.Sprint(v), "Wrong argument[%d] for %s", i, value.Name)
		}
	}
}
