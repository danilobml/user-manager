package errs

import "errors"

var ErrNotFound = errors.New("not found")

var ErrAlreadyExists = errors.New("user with this email already exists")

var ErrParsingRoles = errors.New("user with this email already exists")

var ErrInvalidToken = errors.New("invalid user token")

var ErrParsingToken = errors.New("could not parse user token")

var ErrInvalidCredentials = errors.New("invalid credentials")

var ErrUnauthorized = errors.New("unauthorized")
