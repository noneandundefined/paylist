package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
)

var key []byte
var iv []byte

func Encrypt(plaintext string) (string, error) {
	key, _ = base64.StdEncoding.DecodeString(os.Getenv("SUPER_SECRET_KEY"))
	iv, _ = base64.StdEncoding.DecodeString(os.Getenv("IV"))

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gsm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := gsm.Seal(nil, iv, []byte(plaintext), nil)

	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return encodedCiphertext, nil
}

func Decrypt(encodedCiphertext string) (string, error) {
	key, _ = base64.StdEncoding.DecodeString(os.Getenv("SUPER_SECRET_KEY"))
	iv, _ = base64.StdEncoding.DecodeString(os.Getenv("IV"))

	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gsm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gsm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
