package format

import "fmt"

// FormatUSD formats the given dollars and cents into an expression of US dollars
func FormatUSD(dollars int64, cents int16) string {
	demarcated := fmt.Sprintf("$%d.%02d", abs(dollars), abs(int64(cents)))
	if dollars < 0 || cents < 0 {
		return "-" + demarcated
	}
	return demarcated
}

func abs(v int64) int64 {
	if v < 0 {
		return -1 * v
	}
	return v
}
