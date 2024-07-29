package gns

import (
	"std"
	"time"
)

// IDomain defines the interface for interacting with a domain
type IDomain interface {
	TransferOwnership(newOwner std.Address) error
	Owner() std.Address

	TTL() time.Duration
	ExtendExpiration(duration time.Duration)
	IsAvailable() bool

	Disable()
	Enable()
}
