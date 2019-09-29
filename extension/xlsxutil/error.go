package xlsxutil

import "errors"

var (
	ErrMustBeStruct = errors.New("object must be struct")
	ErrMustHaveSlice = errors.New("object must have slice")
)