package gns

type MetadataField int

const (
	FieldAvatar MetadataField = iota
	FieldRegistrationTime
	FieldRenewalFee
	FieldExpiration
	FieldAttributes
	FieldDescription
	FieldContactInfo
)
