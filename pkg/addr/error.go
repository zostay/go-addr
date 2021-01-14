package addr

type PartialParseError struct {
	Remainder string
}

func (PartialParseError) Error() string {
	return "incomplete parsing of email address"
}
