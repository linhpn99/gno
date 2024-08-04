package gns

import (
	"errors"
)

var (
	ErrCallerIsNotOwner  = errors.New("caller is not token owner")
	ErrInvalidDomainName = errors.New("invalid domain name")
)
