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

	SetDomainData(domainName string, metadata Metadata) error
	GetDomainData(domainName string, field MetadataField) (Metadata, error)
	GetDomainFields(domainName string, fields []MetadataField) (Metadata, error)
}

// NewDomainRegistry creates a new domain registry with metadata extensions
func NewDomainRegistry(name, symbol string) *domainRegistry {
	registry := grc721.NewBasicNFT(name, symbol)

	return &domainRegistry{
		domains:  registry,
		metadata: avl.NewTree(),
	}
}

// SetDomainData sets the metadata for a given domain name
func (s *domainRegistry) SetDomainData(domainName string, metadata Metadata) error {
	owner, err := s.domains.OwnerOf(grc721.TokenID(domainName))
	if err != nil {
		return err
	}

	caller := std.PrevRealm().Addr()
	if caller != owner {
		return ErrCallerIsNotOwner
	}

	s.metadata.Set(domainName, metadata)
	return nil
}

// GetDomainFields retrieves multiple fields of metadata for a given domain
func (s *domainRegistry) GetDomainFields(domainName string, fields []MetadataField) (Metadata, error) {
	data, found := s.metadata.Get(domainName)
	if !found {
		return Metadata{}, ErrInvalidDomainName
	}

	metadata := data.(Metadata)

	if len(fields) == 0 {
		return metadata, nil
	}

	var result Metadata
	for _, field := range fields {
		switch field {
		case FieldAvatar:
			result.Avatar = metadata.Avatar
		case FieldRegistrationTime:
			result.RegistrationTime = metadata.RegistrationTime
		case FieldExpiration:
			result.Expiration = metadata.Expiration
		case FieldAttributes:
			result.Attributes = metadata.Attributes
		case FieldDescription:
			result.Description = metadata.Description
		case FieldContactInfo:
			result.ContactInfo = metadata.ContactInfo
		default:
			return Metadata{}, ErrInvalidMetadataType
		}
	}

	return result, nil
}

// GetDomainData retrieves metadata for a given domain
func (s *domainRegistry) GetDomainData(domainName string, field MetadataField) (Metadata, error) {
	data, found := s.metadata.Get(domainName)
	if !found {
		return Metadata{}, ErrInvalidDomainName
	}

	metadata := data.(Metadata)

	switch field {
	case FieldAvatar:
		return Metadata{
			Avatar: metadata.Avatar,
		}, nil
	case FieldRegistrationTime:
		return Metadata{
			RegistrationTime: metadata.RegistrationTime,
		}, nil
	case FieldExpiration:
		return Metadata{
			Expiration: metadata.Expiration,
		}, nil
	case FieldAttributes:
		return Metadata{
			Attributes: metadata.Attributes,
		}, nil
	case FieldDescription:
		return Metadata{
			Description: metadata.Description,
		}, nil
	case FieldContactInfo:
		return Metadata{
			ContactInfo: metadata.ContactInfo,
		}, nil
	default:
		return Metadata{}, ErrInvalidMetadataType
	}
}
