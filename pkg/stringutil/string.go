package stringutil

import (
	"fmt"
	"math/rand"
	"strconv"
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

func HexStringToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 16, 64)
}

func StringListToInt64(strs []string) []int64 {
	nums := make([]int64, len(strs))
	for i, str := range strs {
		if num, err := strconv.ParseInt(str, 10, 64); err == nil {
			nums[i] = num
		} else {
			return nums[:i]
		}
	}
	return nums
}

func Int64ListToString(list []int64) []string {
	strs := make([]string, len(list))
	for i, num := range list {
		strs[i] = strconv.FormatInt(num, 10)
	}
	return strs
}
