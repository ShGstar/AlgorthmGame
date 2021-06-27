package main

import (
	skip "./skliplist"
)

//test main

func main() {

	slt := skip.NewSkipList()
	for i := 10; i > 0; i-- {
		slt.Insert(i)
		//slt.PrintSkipList()
	}

}
