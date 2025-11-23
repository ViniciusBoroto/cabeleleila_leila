package models

import "errors"

var (
	ErrCannotUpdateWithingTwoDays = errors.New("alterações dentro de 2 dias não são permitidas")
)
