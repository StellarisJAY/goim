package snowflake

import (
	"testing"
)

func BenchmarkSnowflakeID(b *testing.B) {
	snowflake := NewSnowflake(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			snowflake.NextID()
		}
	})
}
