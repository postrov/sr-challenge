package parser

import (
	"fmt"
	"math"
	"testing"

	p "github.com/a-h/parse"
	"github.com/stretchr/testify/assert"
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
		match, ok, err := relativeRowRef.Parse(input)
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
		match, ok, err := floatLiteral.Parse(input)

		assert.True(t, ok)
		assert.Nil(t, err)
		assert.Less(t, math.Abs(c.want-match), epsilon)
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
		match, ok, err := labelRelativeRowRef.Parse(input)
		assert.True(t, ok)
		assert.Nil(t, err)
		assert.Equal(t, c.wantLabelName, match.A)
		assert.Equal(t, c.wantRowNum, match.B)
	}
}

func TestStringLiteralParser(t *testing.T) {
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
		match, ok, err := stringLiteral.Parse(input)
		assert.True(t, ok)
		assert.Nil(t, err)
		assert.Equal(t, c.want, match)
	}
}

/// arExpr

func TestArExprParser(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"15", "15"},
		{"1 + 2 * 3", "1 + (2 * 3)"},
		{"1 + 2 * 3 + 4", "(1 + (2 * 3)) + 4"},
		{"1 + 2 * 3 + 4 * 5 - 6", "((1 + (2 * 3)) + (4 * 5)) - 6"},
	}
	for _, c := range cases {
		input := p.NewInput(c.in)
		match, ok, err := arExprParser.Parse(input)
		assert.True(t, ok)
		assert.Nil(t, err)
		matchRepr := fmt.Sprint(match)
		assert.Equal(t, c.want, matchRepr)
	}
}

func TestIntLitParser(t *testing.T) {
	in := p.NewInput("123")
	m, o, e := parsePrimary.Parse(in)
	assert.True(t, o)
	assert.Nil(t, e)
	assert.Equal(t, IntLit(123), m)

	in = p.NewInput("1 + 2 * 3")
	m, o, e = parsePrimary.Parse(in)
	assert.True(t, o)
	assert.Nil(t, e)
	assert.Equal(t, IntLit(1), m)
}
