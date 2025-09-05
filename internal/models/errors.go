package models

import "errors"

var (
	ErrNoRecord     = errors.New("models: no matching record found")
	ErrTooManyItems = errors.New("too many items in result set")
	ErrTextTooLong  = errors.New("text is too long")
)