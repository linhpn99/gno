package gns

import (
	"std"
	"testing"
	"time"

	"gno.land/p/demo/testutils"
)

func TestGns(t *testing.T) {
	// Test Register
	caller := testutils.TestAddress("foo")
	std.TestSetOrigCaller(caller)

	domainName := "example.gno"
	err := Register(domainName)
	if err != nil {
		t.Fatalf("Register: expected no error, got %v", err)
	}

	// Test Resolve
	address, err := Resolve(domainName)
	if err != nil {
		t.Fatalf("Resolve: expected no error, got %v", err)
	}

	if address != std.PrevRealm().Addr() {
		t.Errorf("Resolve: expected address to be %v, got %v", std.PrevRealm().Addr(), address)
	}

	// Verify that the expiration has been extended
	domain, err := getDomain(domainName)
	if err != nil {
		t.Fatalf("GetDomain: expected no error, got %v", err)
	}

	ttl := domain.TTL()

	// Test ExtendExpiration
	duration := time.Hour * 24 * 30 // Extend by 30 days
	err = ExtendExpiration(domainName, duration)
	if err != nil {
		t.Fatalf("ExtendExpiration: expected no error, got %v", err)
	}

	// expectedExpiration := time.Now().Add(defaultExpireTime + duration)
	// if domain.TTL() > ttl {
	// 	t.Errorf("ExtendExpiration: expected expiration to be extended, got %d", domain.TTL().Unix())
	// }

	// Test Register with invalid name
	invalidDomainName := "invalid name"
	err = Register(invalidDomainName)
	if err == nil {
		t.Fatalf("Register: expected error for invalid domain name, got nil")
	}

	// Test Resolve for non-existent domain
	_, err = Resolve("nonexistent.gno")
	if err == nil {
		t.Fatalf("Resolve: expected error for non-existent domain, got nil")
	}
}
