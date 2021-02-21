package addr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	al, err := ParseEmailAddressList("Some Email <email@example.com>")
	assert.NoError(t, err)
	if !assert.Equal(t, 1, len(al)) {
		return
	}

	if !assert.NotNil(t, al[0]) {
		return
	}

	assert.Equal(t, "email@example.com", al[0].Address())

	ml := al.Flatten()
	if !assert.Equal(t, 1, len(al)) {
		return
	}

	if !assert.NotNil(t, ml[0]) {
		return
	}

	assert.Equal(t, "email@example.com", ml[0].Address())
}
