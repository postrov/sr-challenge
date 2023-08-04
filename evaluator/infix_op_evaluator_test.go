package evaluator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfixOp(t *testing.T) {
	cases := []struct {
		op   func(CalculatedValue, CalculatedValue) CalculatedValue
		lhs  CalculatedValue
		rhs  CalculatedValue
		want CalculatedValue
	}{
		{calcAdd, intValue(3), intValue(4), intValue(7)},
		{calcAdd, stringValue("hello"), stringValue(", world!"), stringValue("hello, world!")},
		{calcAdd, intValue(5), intValue(-3), intValue(2)},
		{calcMul, intValue(5), intValue(-3), intValue(-15)},
		{calcDiv, intValue(5), intValue(2), intValue(2)},
		{calcSub, intValue(100), intValue(-80), intValue(180)},
		{bte_, intValue(10), intValue(7), boolValue(false)},
		{bte_, intValue(5), intValue(5), boolValue(true)},
		{bte_, intValue(4), intValue(5), boolValue(true)},
		{bte_, floatValue(1.0), floatValue(-8.0), boolValue(false)},
		{bte_, floatValue(5.0), intValue(500), boolValue(true)},
		{bte_, intValue(4), floatValue(5), boolValue(true)},
	}

	for _, c := range cases {
		assert.Equal(t, c.want, c.op(c.lhs, c.rhs))
	}
}
