package core

import (
	"encoding/json"
	"errors"
	"fmt"
)

// lineAndColumn calculates 1-based line and column numbers from a byte offset.
func lineAndColumn(data []byte, offset int64) (line, column int, err error) {
	if offset < 0 {
		return 0, 0, fmt.Errorf("invalid offset: %d", offset)
	}
	line = 1
	column = 1
	for i := int64(0); i < int64(len(data)) && i < offset; i++ {
		if data[i] == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}
	return
}

// FormatJSONError adds line and column information to JSON errors when possible.
func FormatJSONError(data []byte, err error, context string) error {
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	switch {
	case errors.As(err, &syntaxErr):
		l, c := lineAndColumn(data, syntaxErr.Offset)
		return fmt.Errorf("%s at line %d, column %d: %w", context, l, c, err)
	case errors.As(err, &typeErr):
		l, c := lineAndColumn(data, typeErr.Offset)
		return fmt.Errorf("%s at line %d, column %d: %w", context, l, c, err)
	default:
		return fmt.Errorf("%s: %w", context, err)
	}
}
