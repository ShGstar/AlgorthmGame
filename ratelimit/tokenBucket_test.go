package ratelimit

import (
	"testing"
	"time"
)

func TestTokenBucket_Access(t *testing.T) {

	tb := NewTokenBucket(time.Second, 10)

	// 前10次操作
	for i := 0; i < 20; i++ {
		if !tb.Access() {
			//fmt.Println("error :",i)
		}
	}

	// 延时1 s
	time.Sleep(time.Second)

	// 第11次操作
	if !tb.Access() {
		//fmt.Println("11 Access fail")
	} else {
		//fmt.Println("11 Access ok")
	}

}
