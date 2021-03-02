// Package encoding defines a custom CharsetReader for expanding the encodings
// that are available for the MIME word decoders used when parsing email
// addresses.
package encoding

import (
	"io"

	_ "golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"

	"github.com/zostay/go-addr/pkg/addr"
)

func init() {
	addr.CharsetReader = CharsetReader
}

// CharsetReader replaces the the CharsetReader of the addr package with one
// that can handle IANA regisered encodings.
func CharsetReader(charset string, r io.Reader) (io.Reader, error) {
	e, err := ianaindex.MIME.Encoding(charset)
	if err != nil {
		return nil, err
	}

	dr := e.NewDecoder().Reader(r)
	return dr, nil
}
