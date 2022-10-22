package data

import (
	"fmt"
	"strconv"
)

// Runtime type, has the same underlying type, int32, same as
// Movie struct field.
type Runtime int32

// MarshalJSON satisfies the json.Marshaler interface.
// Custom implementation to encode the Runtime field on the Movie struct
func (r Runtime) MarshalJSON() ([]byte, error) {
	value := fmt.Sprintf("%d mins", r)

	// Valid JSON is double-quoted
	validJSON := strconv.Quote(value)
	
	return []byte(validJSON), nil
}
