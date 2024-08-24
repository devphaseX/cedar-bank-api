package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var rnd *rand.Rand

const alphabet = "abcdefghijklmnopqrstuvwqyz"

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	rnd = rand.New(source)
}

func RandomInt(min, max int64) int64 {
	return min + rnd.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	s := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rnd.Intn(s)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 10_000)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@mail.com", RandomString(6))
}

func RandomCurrency() string {
	currencies := []string{"USD", "CAD", "EUR"}
	l := len(currencies)

	return currencies[rnd.Intn(l)]
}
