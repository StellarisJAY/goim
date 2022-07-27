package stringutil

import (
	"fmt"
	"math/rand"
	"strings"
)

var chars = []byte("abcdefghijklmnopqrstuvwxyz0123456789")

func RandomString(l int) string {
	s := make([]byte, l)
	for i := 0; i < int(l); i++ {
		s[i] = chars[rand.Intn(len(chars))]
	}
	return string(s)
}

func HexString(bytes []byte) string {
	builder := strings.Builder{}
	for _, b := range bytes {
		b0 := b >> 4
		b1 := b & 0xf
		builder.WriteString(fmt.Sprintf("%x%x", b0, b1))
	}
	return builder.String()
}
