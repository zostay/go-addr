package addr

type Address interface {
	DisplayName() string
	AddressPart() string
	OriginalString() string
	CleanString() string
}

type AddressList []Address
