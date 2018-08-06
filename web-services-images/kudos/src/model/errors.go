package model

import "errors"

var (
	InternalServerError = errors.New("Internal server error")
	NotFoundError       = errors.New("Your requested kudos is not found")
	ConflictError       = errors.New("Your kudos already exist")
)
