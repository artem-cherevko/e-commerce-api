package auth

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNameTaken = errors.New("username already taken")
