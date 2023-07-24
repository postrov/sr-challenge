package parser

import (
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
