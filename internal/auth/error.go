package auth

import "errors"

var PasswordNotMatched error = errors.New("Password not match, you're put wrong password")
var ErrNotFound error = errors.New("User not found")
var ReadError error = errors.New("Cannot read body")
