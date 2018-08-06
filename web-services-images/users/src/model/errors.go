package model

import "errors"

var (
	InternalServerError = errors.New("Internal server error")
	NotFoundError       = errors.New("Your requested user is not found")
	ConflictError       = errors.New("Your user already exist")
)
