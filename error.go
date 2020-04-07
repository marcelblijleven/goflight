package goflight

import "errors"

// ErrInvalidCredentials is returned when a user does not provide a username and/or password for the Goflight client
var ErrInvalidCredentials = errors.New("incorrect client credentials received")

// ErrUnauthorizedAccess is returned when the provided username and password don't have access to the resource
var ErrUnauthorizedAccess = errors.New("you don't have permission to access this resource (403)")

// ErrEndBeforeBegin is returned when the provided end time is before the begin time
var ErrEndBeforeBegin = errors.New("the provided end time is before the begin time")

// ErrTimeRangeTooBig is returned when the time range parameters are too far apart
var ErrTimeRangeTooBig = errors.New("the provided time range is more than the maximum of 2 hours")
