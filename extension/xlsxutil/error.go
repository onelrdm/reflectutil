package xlsxutil

import "errors"

var (
	ErrIsNil = errors.New("object is nil")
	ErrMustBeStruct = errors.New("object must be struct")
	ErrMustHaveSlice = errors.New("object must have slice")
)