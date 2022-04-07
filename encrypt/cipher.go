package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

func newCipherBlock(key string) (cipher.Block, error) {
	hash := sha256.New()
	bkey := []byte(key)
	_, err := hash.Write(bkey)
	if err != nil {
		return nil, err
	}
	cipherKey := hash.Sum(nil)
	return aes.NewCipher(cipherKey)
}

func encryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBEncrypter(block, iv), nil
}

func decryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBDecrypter(block, iv), nil
}

// Encrypt takes in a key and plaintext , encrypts it and
// return a hex representation of the key and plaintext or error.
// This is based on the standard library examples at :
// 	- `https://golang.org/pkg/crypto/cipher/#NewCFBEncrypter`
func Encrypt(key, plaintext string) (string, error) {
	// creating IV for hashing
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream, err := encryptStream(key, iv)
	if err != nil {
		return "", err
	}
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", ciphertext), nil
}

// EncryptWriter returns  a writer that write the encrypted data to the original writer
func EncryptWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream, err := encryptStream(key, iv)
	if err != nil {
		return nil, err
	}
	n, err := w.Write(iv)
	if n != len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to write IV block to writer🖍️")
	}

	return &cipher.StreamWriter{S: stream, W: w}, nil
}

// DecryptReader returns  a reader that decrypt the encrypted data from the original reader
// and provide a way to access it
func DecryptReader(key string, r io.Reader) (*cipher.StreamReader, error) {
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if n < len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to read IV block from reader 👓")
	}
	stream, err := decryptStream(key, iv)
	if err != nil {
		return nil, err
	}
	return &cipher.StreamReader{S: stream, R: r}, nil
}

// Decrypt takes in a key and cipherHex (hex representation of the ciphertext)
// and decrypts it.
// This is based on the standard library examples at :
// 	- `https://golang.org/pkg/crypto/cipher/#NewCFBDecrypter`
func Decrypt(key, cipherHex string) (string, error) {

	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("encrypt: ciphertext too short❕")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream, err := decryptStream(key, iv)
	if err != nil {
		return "", err
	}

	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil

}
