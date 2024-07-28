package gns

import (
	"std"
	"time"

	"gno.land/p/demo/ownable"
)

// Domain struct defines the properties of a domain
type domain struct {
	IDomain

	owner          *ownable.Ownable
	name           string
	creationTime   time.Time
	expirationTime time.Time
	disabled       bool
}

var _ IDomain = (*domain)(nil)

// NewDomain creates a new domain with the given properties
func NewDomain(name string, owner std.Address, expirationTime time.Duration) *domain {
	currentTime := time.Now()

	return &domain{
		owner:          ownable.NewWithAddress(owner),
		name:           name,
		creationTime:   currentTime,
		expirationTime: currentTime.Add(expirationTime),
		disabled:       false,
	}
}

// TransferOwnership transfers the ownership of the domain to a new owner
func (d *domain) TransferOwnership(newOwner std.Address) error {
	return d.owner.TransferOwnership(newOwner)
}

// GetOwner returns the address of the domain owner
func (d domain) Owner() std.Address {
	return d.owner.Owner()
}

// TTL calculates the remaining time to live for the domain
func (d domain) TTL() time.Duration {
	return time.Until(d.expirationTime)
}

// ExtendExpiration extends the expiration time of the domain
func (d *domain) ExtendExpiration(duration time.Duration) {
	d.owner.AssertCallerIsOwner()
	d.expirationTime = d.expirationTime.Add(duration)
}

// IsAvailable checks if the domain is available
func (d domain) IsAvailable() bool {
	return !d.disabled && time.Now().Before(d.expirationTime)
}

// Disable disables the domain
func (d *domain) Disable() {
	d.owner.AssertCallerIsOwner()
	d.disabled = true
}

// Enable enables the domain
func (d *domain) Enable() {
	d.owner.AssertCallerIsOwner()
	d.disabled = false
}
