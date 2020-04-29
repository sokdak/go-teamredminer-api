package cgminer

import "strconv"

// Number is CGMiner API JSON number value that can be
// null, int or empty string.
//
// This type serves as json.Number replacement because
// empty string cannot be unmarshaled into integer
// with json.Number since Go 1.14.
//
// See: https://golang.org/doc/go1.14#encoding/json
type Number string

// String returns the literal text of the number.
func (n Number) String() string { return string(n) }

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error) {
	if n == "" {
		// workaround for empty string literal
		return 0, nil
	}
	return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
	if n == "" {
		// workaround for empty string literal
		return 0, nil
	}
	return strconv.ParseInt(string(n), 10, 64)
}
