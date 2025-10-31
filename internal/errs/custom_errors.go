package errs

import "errors"

var ErrNotFound = errors.New("not found")

var ErrAlreadyExists= errors.New("user with this email already exists")
