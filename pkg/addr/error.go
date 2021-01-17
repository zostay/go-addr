package addr

// PartialParseError is returned when one of the Parse functions is able to
// parse a value out from the start of the string, but was unable to match the
// entire string. This might mean that the string contains additional text after
// the piece it is able to parse or it might mean that the input is formatted
// okay at the start, but contains some unparseable garbage in the middle or
// end. This error allows your implementation to decide whether or not a partial
// parse is acceptable or not.
type PartialParseError struct {
	Remainder string // This is the remaining unparsed string.
}

// Error returns the message "incomplete parsing of email address".
func (PartialParseError) Error() string {
	return "incomplete parsing of email address"
}
