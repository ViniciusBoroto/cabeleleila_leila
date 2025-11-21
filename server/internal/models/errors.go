package models

import "errors"

var (
	ErrCannotUpdateWithingTwoDays = errors.New("cannot update appointment within two days of scheduled date")
)
