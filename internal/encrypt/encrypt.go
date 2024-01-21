package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// Encrypt encrypts the given plain text using the public key located at publicKeyPath.
//
// It takes publicKeyPath and plainText as parameters and returns the encrypted string and an error.
func Encrypt(publicKeyPath, plainText string) (string, error) {
	bytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", err
	}

	publicKey, err := convertBytesToPublicKey(bytes)
	if err != nil {
		return "", err
	}

	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plainText))
	if err != nil {
		return "", err
	}

	return cipherToPemString(cipher), nil
}

// convertBytesToPublicKey converts the given keyBytes to an rsa.PublicKey.
//
// keyBytes []byte - The bytes to be converted to a public key.
// (*rsa.PublicKey, error) - Returns the rsa.PublicKey and an error, if any.
func convertBytesToPublicKey(keyBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes

	cert, err := x509.ParseCertificate(blockBytes)
	if err != nil {
		return nil, err
	}

	return cert.PublicKey.(*rsa.PublicKey), nil
}

// cipherToPemString takes a byte array representing a cipher and returns a string
// representing the PEM encoding of the cipher.
//
// Parameter(s):
//
//	cipher []byte - byte array representing the cipher
//
// Return type(s):
//
//	string - string representing the PEM encoding of the cipher
func cipherToPemString(cipher []byte) string {
	return string(
		pem.EncodeToMemory(
			&pem.Block{
				Type:  "MESSAGE",
				Bytes: cipher,
			},
		),
	)
}

// Decrypt decrypts the encrypted message using the private key.
//
// It takes privateKeyPath and encryptedMessage as parameters and returns a string and an error.
func Decrypt(privateKeyPath, encryptedMessage string) (string, error) {
	bytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	privateKey, err := convertBytesToPrivateKey(bytes)
	if err != nil {
		return "", err
	}

	plainMessage, err := rsa.DecryptPKCS1v15(
		rand.Reader,
		privateKey,
		pemStringToCipher(encryptedMessage),
	)

	return string(plainMessage), nil
}

// convertBytesToPrivateKey converts a byte array to a private key.
//
// keyBytes []byte - the byte array to be converted.
// *rsa.PrivateKey - the converted private key.
// error - an error, if any.
func convertBytesToPrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes

	privateKey, err := x509.ParsePKCS1PrivateKey(blockBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// pemStringToCipher converts a PEM-encoded string to a cipher.
//
// encryptedMessage string
// []byte
func pemStringToCipher(encryptedMessage string) []byte {
	b, _ := pem.Decode([]byte(encryptedMessage))

	return b.Bytes
}
