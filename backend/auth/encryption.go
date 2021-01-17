package auth

import (
	"crypto/aes"
	cryptoCipher "crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type cipher struct {
	gcm cryptoCipher.AEAD
}

func newCipher() *cipher {
	var config struct {
		EncryptionKey string `split_words:"true" required:"true"`
	}
	envconfig.MustProcess("backend", &config)
	c, err := aes.NewCipher([]byte(config.EncryptionKey))
	if err != nil {
		logrus.WithError(err).Panic("invalid secret key provided")
	}
	gcm, err := cryptoCipher.NewGCM(c)
	if err != nil {
		logrus.WithError(err).Panic("invalid cipher generated")
	}
	return &cipher{gcm}
}

func (c *cipher) Encrypt(v string) string {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logrus.WithError(err).Panic("failed to read into nonce")
	}
	return string(c.gcm.Seal(nonce, nonce, []byte(v), nil))
}

func (c *cipher) Decrypt(v string) string {
	vB := []byte(v)
	nonceSize := c.gcm.NonceSize()
	if len(vB) < nonceSize {
		logrus.Panic("trying to decrypt value encrypted by a different cipher")
	}
	nonce, encrypted := vB[:nonceSize], vB[nonceSize:]
	decrypted, err := c.gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		logrus.WithError(err).Panic("failed to open gcm")
	}
	return string(decrypted)
}
