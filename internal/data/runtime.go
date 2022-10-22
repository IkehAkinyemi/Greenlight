package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Runtime type, has the same underlying type, int32, same as
// Movie struct field.
type Runtime int32

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// MarshalJSON satisfies the json.Marshaler interface.
// Custom implementation to encode the Runtime field on the Movie struct
func (r Runtime) MarshalJSON() ([]byte, error) {
	value := fmt.Sprintf("%d mins", r)

	// Valid JSON is double-quoted
	validJSON := strconv.Quote(value)
	
	return []byte(validJSON), nil
}

func (r *Runtime) UnmarshalJSON(json []byte) error {
	// runtime json value is of format "<runtime mins>"
	value, err := strconv.Unquote(string(json))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	parts := strings.Split(value, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	runtime, err := strconv.ParseInt(parts[0], 10, 32)

	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	
	*r = Runtime(runtime)

	return nil
}