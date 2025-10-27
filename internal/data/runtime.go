package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	if r < 0 {
		return nil, ErrInvalidRuntimeFormat
	}

	var jsonValue string
	if r == 1 {
		jsonValue = fmt.Sprintf("%d min", r)
	} else {
		jsonValue = fmt.Sprintf("%d mins", r)
	}

	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 {
		return ErrInvalidRuntimeFormat
	}

	inte, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil || inte < 0 {
		return ErrInvalidRuntimeFormat
	}

	expectedUnit := "mins"
	if inte == 1 {
		expectedUnit = "min"
	}

	if parts[1] != expectedUnit {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(inte)
	return nil
}
