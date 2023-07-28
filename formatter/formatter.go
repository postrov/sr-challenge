package formatter

import "strconv"

// Common float formatter
func Ftoa(f float64) string {
	return strconv.FormatFloat(f, 'f', 3, 64)
}
