package http

// ParseMillidollars parses the "millidollar" unit returned by YNAB into dollars and cents
func ParseMillidollars(millidollars int64) (int64, int16) {
	asCents := millidollars % 1000
	asDollars := (millidollars - asCents) / 1000
	return asDollars, int16(asCents) / 10
}

// ToMillidollars converts the given dollars and cents to millidollars for YNAB's API
func ToMillidollars(dollars int64, cents int16) int64 {
	return dollars*1000 + int64(cents)*10
}
