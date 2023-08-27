package keyring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyring(t *testing.T) {
	kr := NewKeyring("test", &Config{
		Backend: "sagecipher.keyring.Keyring",
	})
	theSecret := "cookiemonsta"
	err := kr.Set("secret", theSecret)
	assert.Nil(t, err)
	secret, err := kr.Get("secret")
	assert.Nil(t, err)
	assert.Equal(t, theSecret, secret)
	err = kr.Del("secret")
	assert.Nil(t, err)
}
