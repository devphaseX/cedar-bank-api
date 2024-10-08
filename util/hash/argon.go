package hash

import (
	"bytes"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

type Argon2idHash struct {
	// time represents the number of
	// passed over the specified memory.
	time uint32
	// cpu memory to be used.
	memory uint32
	// threads for parallelism aspect
	// of the algorithm.
	threads uint8
	// keyLen of the generate hash key.
	keyLen uint32
	// saltLen the length of the salt used.
	saltLen uint32
}

// NewArgon2idHash constructor function for
// Argon2idHash.
func NewArgon2idHash(time, saltLen uint32, memory uint32, threads uint8, keyLen uint32) *Argon2idHash {
	return &Argon2idHash{
		time:    time,
		saltLen: saltLen,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
	}
}

type HashSalt struct {
	Hash, Salt []byte
}

func (hs *HashSalt) EncodeBase64() (string, string) {
	return ArgonStringEncode(hs)
}

func ArgonStringEncode(hs *HashSalt) (string, string) {
	hashedPasswordStr := base64.StdEncoding.EncodeToString(hs.Hash)
	saltStr := base64.StdEncoding.EncodeToString(hs.Salt)
	return hashedPasswordStr, saltStr
}

func ArgonStringDecode(passwordHashStr, passwordSaltStr string) ([]byte, []byte) {
	hashedPassword, _ := base64.StdEncoding.DecodeString(passwordHashStr)
	salt, _ := base64.StdEncoding.DecodeString(passwordSaltStr)
	return hashedPassword, salt
}

// GenerateHash using the password and provided salt.
// If not salt value provided fallback to random value
// generated of a given length.
func (a *Argon2idHash) GenerateHash(password, salt []byte) (*HashSalt, error) {
	var err error
	// If salt is not provided generate a salt of
	// the configured salt length.
	if len(salt) == 0 {
		salt, err = randomSecret(a.saltLen)
	}
	if err != nil {
		return nil, err
	}
	// Generate hash
	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keyLen)
	// Return the generated hash and salt used for storage.
	return &HashSalt{Hash: hash, Salt: salt}, nil
}

// Compare generated hash with store hash.
func (a *Argon2idHash) Compare(hash, salt, password []byte) error {
	// Generate hash for comparison.
	hashSalt, err := a.GenerateHash(password, salt)
	if err != nil {
		return err
	}
	// Compare the generated hash with the stored hash.
	// If they don't match return error.
	if !bytes.Equal(hash, hashSalt.Hash) {
		return errors.New("hash doesn't match")
	}
	return nil
}

type ArgonHashPayload struct {
	PasswordHash string
	PasswordSalt string
}

func DefaultArgonHash() *Argon2idHash {
	return NewArgon2idHash(1, 32, 64*1024, 32, 256)
}
