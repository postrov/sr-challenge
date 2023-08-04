package evaluator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"pasza.org/sr-challenge/model"
)

func TestIncFrom(t *testing.T) {
	csvCells := make(CSVCells, 1)
	csvCells[0] = make([]model.Cell, 1)
	csvCells[0][0] = model.IntCell{
		Value: 0,
	}
	evalCells := make([][]evalCell, 1)
	evalCells[0] = make([]evalCell, 1)
	evalCells[0][0] = evalCell{
		done:      false,
		copyCount: 5,
		formula:   nil,
		value:     intValue(8),
	}
	ec := &evalCells[0][0]
	lm := make(labelMap)
	es := evalState{
		evalCells:   evalCells,
		csvCells:    csvCells,
		labelsOnRow: []labelMap{lm},
	}
	cases := []struct {
		testName   string
		fn         func(*evalState, []CalculatedValue, int, int) CalculatedValue
		args       []CalculatedValue
		want       CalculatedValue
		shouldFail bool
	}{
		{
			"incFrom should fail with no args",
			incFrom,
			[]CalculatedValue{},
			intValue(0),
			true,
		},
		{
			"incFrom should return copyCount + args[0] value",
			incFrom,
			[]CalculatedValue{intValue(1)},
			intValue(ec.copyCount + 1),
			false,
		},
	}

	for _, c := range cases {
		if c.shouldFail {
			assert.Panics(t, func() { c.fn(&es, c.args, 0, 0) })
		} else {
			actual := c.fn(&es, c.args, 0, 0)
			assert.Equal(t, c.want, actual)
		}
	}
}
