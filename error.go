package goflight

import "errors"

var InvalidCredentialsError = errors.New("incorrect client credentials received")
var UnauthorizedAccessError = errors.New("you don't have permission to access this resource (403)")
