package parser

import m "pasza.org/sr-challenge/model"

// Simplified arithmetic operator precedence fix for an expression, based on pseudocode from:
// https://en.wikipedia.org/wiki/Operator-precedence_parser#Pseudocode

// simplified operator priority fix
func fixOperatorPrecedence(primaries []ArExpr, binOps []m.BinOp) ArExpr {

	lhs := primaries[0]
	res, _, _ := fixOperatorPrecedenceRec(lhs, primaries, binOps, 1, 0, 0)
	return res
}

func fixOperatorPrecedenceRec(lhs ArExpr, primaries []ArExpr, binOps []m.BinOp, pIndex int, oIndex int, minPrecedence int) (res ArExpr, newPIndex int, newOIndex int) {
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
		lhs = Op{
			lhs: lhs,
			rhs: rhs,
			op:  op,
		}
	}
	return lhs, pIndex, oIndex
}
