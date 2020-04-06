package goflight

import "errors"

// ErrInvalidCredentials is returned when a user does not provide a username and/or password for the Goflight Client
var ErrInvalidCredentials = errors.New("incorrect Client credentials received")
// ErrUnauthorizedAccess is returned when the provided username and password don't have access to the resource
var ErrUnauthorizedAccess = errors.New("you don't have permission to access this resource (403)")

