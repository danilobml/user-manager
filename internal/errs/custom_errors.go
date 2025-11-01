package errs

import "errors"

var ErrNotFound = errors.New("not found")

var ErrAlreadyExists = errors.New("user with this email already exists")

var ErrParsingRoles = errors.New("user with this email already exists")

var ErrInvalidToken = errors.New("invalid user token")

var ErrInvalidCredentials = errors.New("invalid credentials")
