package gns

import (
	"std"
	"testing"
	"time"

	"gno.land/p/demo/ownable"
	"gno.land/p/demo/testutils"
	"gno.land/p/demo/uassert"
)

var (
	addr1 = testutils.TestAddress("foo")
	addr2 = testutils.TestAddress("bar")
)

func TestNewDomain(t *testing.T) {
	owner := addr1
	name := "example.gno"
	expirationTime := time.Hour * 24 * 30

	domain := NewDomain(name, owner, expirationTime)

	uassert.Equal(t, domain.Owner(), owner)
	uassert.Equal(t, domain.name, name)
	uassert.Equal(t, domain.disabled, false)
	// uassert.WithinDuration(t, domain.expirationTime, time.Now().Add(expirationTime), time.Second)
}

func TestTransferOwnership(t *testing.T) {
	owner := addr1
	newOwner := addr2

	domain := NewDomain("example.gno", owner, time.Hour*24*30)

	std.TestSetRealm(std.NewUserRealm(addr1))
	std.TestSetOrigCaller(addr1) // TODO(bug): should not be needed

	err := domain.TransferOwnership(newOwner)
	uassert.NoError(t, err)
	uassert.Equal(t, domain.Owner(), newOwner)
}

func TestIsAvailable(t *testing.T) {
	owner := addr1
	domain := NewDomain("example.gno", owner, time.Hour*24*30)

	uassert.Equal(t, domain.IsAvailable(), true)

	domain.Disable()
	uassert.Equal(t, domain.IsAvailable(), false)

	domain.Enable()
	uassert.Equal(t, domain.IsAvailable(), true)

	domain.expirationTime = time.Now().Add(-time.Hour)
	uassert.Equal(t, domain.IsAvailable(), false)
}

func TestTTL(t *testing.T) {
	owner := addr1
	expirationTime := time.Hour * 24 * 30
	domain := NewDomain("example.gno", owner, expirationTime)

	ttl := domain.TTL()
	// uassert.WithinDuration(t, ttl, expirationTime, time.Second)
}

func TestExtendExpiration(t *testing.T) {
	owner := addr1
	expirationTime := time.Hour * 24 * 30
	domain := NewDomain("example.gno", owner, expirationTime)

	domain.ExtendExpiration(time.Hour * 24)
	// uassert.WithinDuration(t, domain.expirationTime, time.Now().Add(expirationTime+time.Hour*24), time.Second)
}

func TestDisableEnable(t *testing.T) {
	owner := addr1
	domain := NewDomain("example.gno", owner, time.Hour*24*30)

	domain.Disable()
	uassert.Equal(t, domain.disabled, true)

	domain.Enable()
	uassert.Equal(t, domain.disabled, false)
}
