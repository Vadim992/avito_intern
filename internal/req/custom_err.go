package req

import "errors"

var (
	QueryDataErr  = errors.New("incorrect data in query")
	EmptyFieldErr = errors.New("dto struct has empty field")

	OldDataErr = errors.New("data in InMemory is old")

	PathIdErr = errors.New("incorrect id in path")
)
