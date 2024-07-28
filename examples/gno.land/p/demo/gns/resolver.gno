package gns

import (
	"errors"
	"std"

	"gno.land/p/demo/avl"
)

type IResolver interface {
	SetAddress(name string, address std.Address)
	Resolve(name string) (std.Address, error)
}

type resolver struct {
	addresses *avl.Tree
}

func NewResolver() *resolver {
	return &resolver{
		addresses: avl.NewTree(),
	}
}

func (r *resolver) SetAddress(name string, address std.Address) {
	r.addresses.Set(name, address)
}

func (r *resolver) Resolve(name string) (std.Address, error) {
	value, exists := r.addresses.Get(name)
	if !exists {
		return "", errors.New("domain does not exist")
	}
	return value.(std.Address), nil
}
