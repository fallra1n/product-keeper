package shared

import "errors"

var (
	ErrNoData   = errors.New("data not found")
	ErrInternal = errors.New("internal error, please try again later")
)
