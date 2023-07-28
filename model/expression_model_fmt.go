package model

// implement Stringer for expressions for nice display

import (
	"fmt"
	"strconv"
	"strings"

	"pasza.org/sr-challenge/formatter"
)

func (op BinaryOperator) String() string {
	switch op {
	case MUL:
		return "*"
	case DIV:
		return "/"
	case ADD:
		return "+"
	case SUB:
		return "-"
	default:
		return "?"
	}
}

func formatInfixOperand(e Expr) string {
	switch e.(type) {
	case InfixOp:
		return fmt.Sprintf("(%v)", e)
	default:
		return fmt.Sprint(e)
	}
}

func (op InfixOp) String() string {
	lhs, rhs := formatInfixOperand(op.Lhs), formatInfixOperand(op.Rhs)
	return fmt.Sprintf("%v %v %v", lhs, op.Op, rhs)
}

func (cellRef CellRef) String() string {
	return fmt.Sprintf("%s%d", cellRef.Col, cellRef.Row)
}

func (CopyAbove) String() string {
	return "^^"
}

func (v CopyLastInColumn) String() string {
	return fmt.Sprintf("%s^v", v.Col)
}

func (v CopyColumnAbove) String() string {
	return fmt.Sprintf("%s^", v.Col)
}

func (v IntLit) String() string {
	return strconv.Itoa(int(v))
}

func (v FloatLit) String() string {
	return formatter.Ftoa(float64(v))
}

func (s StringLit) String() string {
	escaped := strings.ReplaceAll(string(s), `"`, `\"`)
	return fmt.Sprintf("\"%s\"", escaped)
}

func (fc FunCall) String() string {
	params := make([]string, len(fc.Params))
	for i, expr := range fc.Params {
		params[i] = fmt.Sprint(expr)
	}
	return fmt.Sprintf("%s(%s)", fc.Name, strings.Join(params, ", "))
}
