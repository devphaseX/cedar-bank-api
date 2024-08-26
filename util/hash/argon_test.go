package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	argonHash := DefaultArgonHash()
	password := "verysecurepassword"

	passwordHash, err := argonHash.GenerateHash([]byte(password), nil)

	require.NoError(t, err)

	passwordHashStr, passwordSaltStr := ArgonStringEncode(passwordHash)

	passwordHashByte, passwordSaltByte := ArgonStringDecode(passwordHashStr, passwordSaltStr)
	err = argonHash.Compare(passwordHashByte, passwordSaltByte, []byte(password))
	require.NoError(t, err)
}

func TestIncorrectPassword(t *testing.T) {
	argonHash := DefaultArgonHash()
	password := "verysecurepassword"

	passwordHash, err := argonHash.GenerateHash([]byte(password), nil)

	require.NoError(t, err)

	passwordHashStr, passwordSaltStr := ArgonStringEncode(passwordHash)

	passwordHashByte, passwordSaltByte := ArgonStringDecode(passwordHashStr, passwordSaltStr)

	enteredPassword := "wrongpassword"
	err = argonHash.Compare(passwordHashByte, passwordSaltByte, []byte(enteredPassword))
	require.Error(t, err)
}
