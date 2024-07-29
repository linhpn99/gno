package gns

import (
	"errors"
	"regexp"
	"std"
	"time"

	"gno.land/p/demo/avl"
	"gno.land/p/demo/gns"
)

var (
	domains           *avl.Tree // domain => IDomain
	resolver          gns.IResolver
	defaultExpireTime = time.Hour * 24 * 30 // 30 days
	reName            = regexp.MustCompile(`^[a-zA-Z0-9]{1,124}\.gno$`)
)

func init() {
	domains = avl.NewTree()
	resolver = gns.NewResolver()
}

func isValidName(name string) bool {
	return reName.MatchString(name)
}

func Register(name string, address std.Address) {
	if !isValidName(name) {
		panic("invalid domain name")
	}

	if domains.Has(name) {
		panic("domain name has been registered before")
	}

	caller := std.PrevRealm().Addr()
	domain := gns.NewDomain(name, caller, defaultExpireTime)

	domains.Set(name, domain)
	resolver.SetAddress(name, address)
}

func Resolve(name string) std.Address {
	domain, err := getDomain(name)
	if err != nil {
		panic(err)
	}

	if !domain.IsAvailable() {
		panic("domain is unavailable")
	}

	addr, err := resolver.Resolve(name)
	if err != nil {
		panic(err)
	}

	return addr
}

func ExtendExpiration(name string, duration time.Duration) {
	domain, err := getDomain(name)
	if err != nil {
		panic(err)
	}

	domain.ExtendExpiration(duration)
}

func getDomain(name string) (gns.IDomain, error) {
	value, exists := domains.Get(name)
	if !exists {
		return nil, errors.New("domain does not exist")
	}

	return value.(gns.IDomain), nil
}
