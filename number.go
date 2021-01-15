package cgminer

import (
	"fmt"
	"math"
	"strconv"
)

// Number is CGMiner API JSON number value that can be
// null, int or empty string.
//
// This type serves as json.Number replacement because
// empty string cannot be unmarshaled into integer
// with json.Number since Go 1.14.
//
// See: https://golang.org/doc/go1.14#encoding/json
type Number float64

// String returns the literal text of the number.
func (n Number) String() string {
	return strconv.FormatFloat(float64(n), 'E', -1, 64)
}

// Float64 returns the number as a float64.
func (n Number) Float64() float64 {
	return float64(n)
}

// Int64 returns the number as an int64.
func (n Number) Int64() int64 {
	return int64(n.Int())
}

// Int returns the number as an int.
func (n Number) Int() int {
	return int(math.Round(float64(n)))
}

func (n *Number) UnmarshalJSON(b []byte) error {
	str := string(b)
	switch str {
	case "", `""`, "null":
		return nil
	default:
	}

	switch str[0] {
	case '"':
		unquoted, err := strconv.Unquote(str)
		if err != nil {
			return err
		}
		str = unquoted
	case '{', '[':
		return fmt.Errorf("Number.UnmarshalJSON: value is not a number  - %s", str)
	}

	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	*n = Number(num)
	return nil
}

func (n *Number) MarshalJSON() ([]byte, error) {
	return []byte(n.String()), nil
}
