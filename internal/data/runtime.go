package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	var jsonValue string
	if r == 1 {
		jsonValue = fmt.Sprintf("%d min", r)
	} else {
		jsonValue = fmt.Sprintf("%d mins", r)
	}

	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}
