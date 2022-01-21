package encryption

import (
	"strconv"
	"testing"
)

func BenchmarkGetUniqueId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//fmt.Println(GetUniqueId(strconv.Itoa(i)))
		GetUniqueId(strconv.Itoa(i))
	}
}
