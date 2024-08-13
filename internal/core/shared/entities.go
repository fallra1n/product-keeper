package shared

import "errors"

var (
	// ErrNoData data not found
	ErrNoData = errors.New("data not found")

	// ErrInternal internal error, please try again later
	ErrInternal = errors.New("internal error, please try again later")
)
