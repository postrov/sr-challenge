package evaluator

import (
	"strconv"

	"pasza.org/sr-challenge/formatter"
)

// boring infix type checks here
func calcMul(lhs, rhs CalculatedValue) CalculatedValue {
	switch l := lhs.(type) {
	case intValue:
		switch r := rhs.(type) {
		case intValue:
			return intValue(l * r)
		case floatValue:
			return floatValue(float64(l) * float64(r))
		default:
			panic("Unknown value (int * ?)")

		}
	case floatValue:
		switch r := rhs.(type) {
		case intValue:
			return floatValue(float64(l) * float64(r))
		case floatValue:
			return floatValue(l * r)
		default:
			panic("Unknown value (float * ?)")

		}
	case stringValue:
		panic("Multiplication not supported for strings")
	default:
		panic("Unknown value (? * ?)")

	}
}

// fixme: div by 0
func calcDiv(lhs, rhs CalculatedValue) CalculatedValue {
	switch l := lhs.(type) {
	case intValue:
		switch r := rhs.(type) {
		case intValue:
			return intValue(l / r)
		case floatValue:
			return floatValue(float64(l) / float64(r))
		default:
			panic("Unknown value (int / ?)")

		}
	case floatValue:
		switch r := rhs.(type) {
		case intValue:
			return floatValue(float64(l) / float64(r))
		case floatValue:
			return floatValue(l / r)
		default:
			panic("Unknown value (float / ?)")

		}
	case stringValue:
		panic("Division not supported for strings")
	default:
		panic("Unknown value (? / ?)")

	}
}

func calcAdd(lhs, rhs CalculatedValue) CalculatedValue {
	switch l := lhs.(type) {
	case intValue:
		switch r := rhs.(type) {
		case intValue:
			return intValue(l + r)
		case floatValue:
			return floatValue(float64(l) + float64(r))
		default:
			panic("Unknown value (int + ?)")

		}
	case floatValue:
		switch r := rhs.(type) {
		case intValue:
			return floatValue(float64(l) + float64(r))
		case floatValue:
			return floatValue(l + r)
		default:
			panic("Unknown value (float + ?)")

		}
	case stringValue:
		switch r := rhs.(type) {
		case intValue:
			return stringValue(l + stringValue(strconv.Itoa(int(r))))
		case floatValue:
			return stringValue(l + stringValue(formatter.Ftoa(float64(r))))
		case stringValue:
			return stringValue(l + r)
		default:
			panic("Unknown value (string + ?)")
		}
	default:
		panic("Unknown value (? + ?)")
	}
}

func calcSub(lhs, rhs CalculatedValue) CalculatedValue {
	switch l := lhs.(type) {
	case intValue:
		switch r := rhs.(type) {
		case intValue:
			return intValue(l - r)
		case floatValue:
			return floatValue(float64(l) - float64(r))
		default:
			panic("Unknown value (int - ?)")

		}
	case floatValue:
		switch r := rhs.(type) {
		case intValue:
			return floatValue(float64(l) - float64(r))
		case floatValue:
			return floatValue(l - r)
		default:
			panic("Unknown value (float - ?)")

		}
	case stringValue:
		panic("Subtraction not supported for strings")
	default:
		panic("Unknown value (? - ?)")
	}
}

func bte_(lhs, rhs CalculatedValue) CalculatedValue {
	switch l := lhs.(type) {
	case intValue:
		switch r := rhs.(type) {
		case intValue:
			return boolValue(l <= r)
		case floatValue:
			return boolValue(float64(l) <= float64(r))
		}
	case floatValue:
		switch r := rhs.(type) {
		case intValue:
			return boolValue(float64(l) <= float64(r))
		case floatValue:
			return boolValue(l <= r)
		}
	}
	panic("Unsupported operand types for BTE")
}
