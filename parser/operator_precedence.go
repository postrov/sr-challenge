package parser

import m "pasza.org/sr-challenge/model"

// Simplified arithmetic operator precedence fix for an expression, based on pseudocode from:
// https://en.wikipedia.org/wiki/Operator-precedence_parser#Pseudocode

var opPrecedence map[m.BinaryOperator]int = map[m.BinaryOperator]int{
	m.MUL: 2,
	m.DIV: 2,
	m.ADD: 1,
	m.SUB: 1,
}

// simplified operator priority fix
func fixOperatorPrecedence(primaries []m.Expr, binOps []m.BinaryOperator) m.Expr {
	lhs := primaries[0]
	res, _, _ := fixOperatorPrecedenceRec(lhs, primaries, binOps, 1, 0, 0)
	return res
}

func fixOperatorPrecedenceRec(lhs m.Expr, primaries []m.Expr, binOps []m.BinaryOperator, pIndex int, oIndex int, minPrecedence int) (res m.Expr, newPIndex int, newOIndex int) {
	if pIndex >= len(primaries) {
		return lhs, pIndex, pIndex - 1
	}
	lookahead := binOps[oIndex]
	for opPrecedence[lookahead] >= minPrecedence {
		op := lookahead
		oIndex++
		opPrec := opPrecedence[op]
		if pIndex >= len(primaries) {
			break
		}
		rhs := primaries[pIndex]
		pIndex++

		for oIndex < len(binOps) {
			lookahead = binOps[oIndex]
			if opPrecedence[lookahead] <= opPrec {
				break
			}
			rhs, pIndex, oIndex = fixOperatorPrecedenceRec(rhs, primaries, binOps, pIndex, oIndex, opPrec+1)
		}
		lhs = m.InfixOp{
			Lhs: lhs,
			Rhs: rhs,
			Op:  op,
		}
	}
	return lhs, pIndex, oIndex
}
