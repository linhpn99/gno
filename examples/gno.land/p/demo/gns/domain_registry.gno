package gns

import (
	"std"
	"time"

	"gno.land/p/demo/avl"
	"gno.land/p/demo/grc/grc721"
	"gno.land/p/demo/ownable"
)

// domainRegistry represents a registry for domain names with metadata
type domainRegistry struct {
	domains  grc721.IGRC721 // Interface for basic NFT functionality
	metadata *avl.Tree      // AVL tree for storing domain metadata
}

// DomainRegistry defines the methods for managing domain names and metadata
type DomainRegistry interface {
	grc721.IGRC721

	SetMetadata(name DomainName, metadata Metadata) error
	GetMetadata(name DomainName) (Metadata, error)
}

// NewDomainRegistry creates a new domain registry with metadata extensions
func NewDomainRegistry(name, symbol string) *domainRegistry {
	registry := grc721.NewBasicNFT(name, symbol)

	return &domainRegistry{
		domains:  registry,
		metadata: avl.NewTree(),
	}
}

// SetMetadata sets metadata for a given domain
func (s *domainRegistry) SetMetadata(name DomainName, metadata Metadata) error {
	owner, err := s.domains.OwnerOf(grc721.TokenID(name))
	if err != nil {
		return err
	}

	caller := std.PrevRealm().Addr()
	if caller != owner {
		return ErrCallerIsNotOwner
	}

	s.metadata.Set(string(name), metadata)
	return nil
}

// GetMetadata retrieves metadata for a given domain
func (s *domainRegistry) GetMetadata(name DomainName) (Metadata, error) {
	metadata, found := s.metadata.Get(string(name))
	if !found {
		return Metadata{}, ErrInvalidDomainName
	}

	return metadata.(Metadata), nil
}
