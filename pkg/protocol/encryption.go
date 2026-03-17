package protocol

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// Encryptor handles packet encryption/decryption using AES-CFB8.
//
// Minecraft uses AES encryption with CFB8 mode for network packets.
type Encryptor struct {
	key []byte
	// FIXED: CFB8 mode requires separate IV for encryption and decryption
	encryptIV []byte
	decryptIV []byte
	encryptor cipher.Stream
	decryptor cipher.Stream
}

// NewEncryptor creates a new encryptor with the given shared secret.
// The key should be 16 bytes (128-bit) for AES-128.
func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) != 16 {
		return nil, errors.New("encryption key must be 16 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// FIXED: In CFB8 mode, the IV is typically sent separately or derived
	// For Minecraft, we use the first block received from server as IV
	// Initialize with zeros, will be set when encryption is enabled
	encryptIV := make([]byte, aes.BlockSize)
	decryptIV := make([]byte, aes.BlockSize)

	return &Encryptor{
		key:        key,
		encryptIV:  encryptIV,
		decryptIV:  decryptIV,
		encryptor:  cipher.NewCFB8Encrypter(block, encryptIV),
		decryptor:  cipher.NewCFB8Decrypter(block, decryptIV),
	}, nil
}

// Key returns the encryption key.
func (e *Encryptor) Key() []byte {
	return e.key
}

// Encrypt encrypts data using AES-CFB8.
func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
	if len(e.key) != 16 {
		return nil, errors.New("invalid encryption key length")
	}

	encrypted := make([]byte, len(data))
	e.encryptor.XORKeyStream(encrypted, data)

	return encrypted, nil
}

// Decrypt decrypts data using AES-CFB8.
func (e *Encryptor) Decrypt(data []byte) ([]byte, error) {
	if len(e.key) != 16 {
		return nil, errors.New("invalid decryption key length")
	}

	decrypted := make([]byte, len(data))
	e.decryptor.XORKeyStream(decrypted, data)

	return decrypted, nil
}

// SetEncryptionIV sets the IV for encryption/decryption.
// This is called when receiving the encryption request from server.
func (e *Encryptor) SetEncryptionIV(iv []byte) error {
	if len(iv) != aes.BlockSize {
		return errors.New("IV must be 16 bytes")
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	e.encryptIV = iv
	e.decryptIV = iv
	e.encryptor = cipher.NewCFB8Encrypter(block, iv)
	e.decryptor = cipher.NewCFB8Decrypter(block, iv)

	return nil
}

// EncryptWriter wraps an io.Writer and encrypts data written to it.
type EncryptWriter struct {
	writer io.Writer
	stream cipher.Stream
}

// NewEncryptWriter creates a new encrypt writer.
func NewEncryptWriter(writer io.Writer, key []byte) (*EncryptWriter, error) {
	if len(key) != 16 {
		return nil, errors.New("encryption key must be 16 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// FIXED: For stream encryption, IV should be unique per encryption
	// In this case, use zeros as IV (will be set by SetEncryptionIV)
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCFB8Encrypter(block, iv)

	return &EncryptWriter{
		writer: writer,
		stream: stream,
	}, nil
}

// Write encrypts and writes data to the underlying writer.
func (ew *EncryptWriter) Write(data []byte) (int, error) {
	encrypted := make([]byte, len(data))
	ew.stream.XORKeyStream(encrypted, data)
	return ew.writer.Write(encrypted)
}

// Close closes the encrypt writer.
func (ew *EncryptWriter) Close() error {
	// Nothing to close for stream cipher
	return nil
}

// DecryptReader wraps an io.Reader and decrypts data read from it.
type DecryptReader struct {
	reader io.Reader
	stream cipher.Stream
}

// NewDecryptReader creates a new decrypt reader.
func NewDecryptReader(reader io.Reader, key []byte) (*DecryptReader, error) {
	if len(key) != 16 {
		return nil, errors.New("decryption key must be 16 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// FIXED: Use zeros as IV initially
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCFB8Decrypter(block, iv)

	return &DecryptReader{
		reader: reader,
		stream: stream,
	}, nil
}

// Read reads and decrypts data from the underlying reader.
func (dr *DecryptReader) Read(data []byte) (int, error) {
	encrypted := make([]byte, len(data))
	n, err := io.ReadFull(dr.reader, encrypted)
	if err != nil && err != io.ErrUnexpectedEOF {
		return 0, err
	}

	dr.stream.XORKeyStream(data[:n], encrypted[:n])
	return n, err
}

// Close closes the decrypt reader.
func (dr *DecryptReader) Close() error {
	return nil
}

// GenerateSharedSecret generates a random 16-byte shared secret for encryption.
// This is used during the login process to generate the encryption key.
func GenerateSharedSecret() ([]byte, error) {
	secret := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, secret); err != nil {
		return nil, fmt.Errorf("failed to generate shared secret: %w", err)
	}
	return secret, nil
}

// GenerateVerificationToken generates a random 4-byte verification token.
// Used in the encryption handshake.
func GenerateVerificationToken() ([]byte, error) {
	token := make([]byte, 4)
	if _, err := io.ReadFull(rand.Reader, token); err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}
	return token, nil
}

// GenerateNonce generates a random nonce for cryptographic operations.
func GenerateNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	return nonce, nil
}

// bytesEqual securely compares two byte slices in constant time.
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// EncryptAESCBC encrypts data using AES-CBC mode with PKCS#7 padding.
// Not used in standard Minecraft protocol but provided for completeness.
func EncryptAESCBC(data, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Apply PKCS#7 padding
	data = pkcs7Pad(data, aes.BlockSize)

	ciphertext := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, data)

	return ciphertext, nil
}

// DecryptAESCBC decrypts data using AES-CBC mode with PKCS#7 padding.
func DecryptAESCBC(data, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	plaintext := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, data)

	// Remove PKCS#7 padding
	plaintext, err = pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// pkcs7Pad applies PKCS#7 padding to data.
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs7Unpad removes PKCS#7 padding from data.
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	padding := int(data[len(data)-1])
	if padding < 1 || padding > blockSize {
		return nil, errors.New("invalid padding")
	}

	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			return nil, errors.New("invalid padding")
		}
	}

	return data[:len(data)-padding], nil
}
