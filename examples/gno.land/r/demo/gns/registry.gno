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

func Register(name string) error {
	if !isValidName(name) {
		return errors.New("invalid domain name")
	}

	caller := std.PrevRealm().Addr()
	domain := gns.NewDomain(name, caller, defaultExpireTime)

	domains.Set(name, domain)
	resolver.SetAddress(name, domain.Owner())

	return nil
}

func getDomain(name string) (gns.IDomain, error) {
	value, exists := domains.Get(name)
	if !exists {
		return nil, errors.New("domain does not exist")
	}
	return value.(gns.IDomain), nil
}

func Resolve(name string) (std.Address, error) {
	domain, err := getDomain(name)
	if err != nil {
		return "", err
	}
	if !domain.IsAvailable() {
		return "", errors.New("domain is unavailable")
	}
	return resolver.Resolve(name)
}

func ExtendExpiration(name string, duration time.Duration) error {
	domain, err := getDomain(name)
	if err != nil {
		return err
	}
	domain.ExtendExpiration(duration)
	return nil
}
