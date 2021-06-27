package skliplist

import (
	"fmt"
	"math/rand"
)

/*
	结构体创建value为interface{}
	这里其实value为int的跳跃表
*/

//最大层数
const SKIPLIST_MAX_LEVEL = 8

type Node struct {
	Forward []Node
	Value   interface{}
}

func NewNode(v interface{}, level int) *Node {
	return &Node{
		Value:   v,
		Forward: make([]Node, level),
	}
}

type SkipList struct {
	Header *Node
	Level  int
}

func NewSkipList() *SkipList {
	return &SkipList{
		Header: NewNode(0, SKIPLIST_MAX_LEVEL),
		Level:  1,
	}
}

//插入的int值
func (skipList *SkipList) Insert(key int) {
	update := make(map[int]*Node)
	node := skipList.Header

	for i := skipList.Level - 1; i >= 0; i-- {
		for {
			if node.Forward[i].Value != nil && node.Forward[i].Value.(int) < key {
				node = &node.Forward[i]
			} else {
				break
			}
		}
		update[i] = node
	}

	level := skipList.RandomLevel()

	if level > skipList.Level {
		for i := skipList.Level; i < level; i++ {
			update[i] = skipList.Header
		}
		skipList.Level = level
	}

	newNode := NewNode(key, level)

	for i := 0; i < level; i++ {
		newNode.Forward[i] = update[i].Forward[i]
		update[i].Forward[i] = *newNode
	}

}

func (skipList *SkipList) Search(key int) (*Node, bool) {
	node := skipList.Header
	for i := skipList.Level - 1; i >= 0; i++ {
		if node.Forward[i].Value == nil {
			break
		}

		if node.Forward[i].Value.(int) == key {
			return &node.Forward[i], true
		}

		if node.Forward[i].Value.(int) < key { //目标值大于当前的
			node = &node.Forward[i]
			continue
		} else {
			break
		}
	}

	return nil, false
}

func (skipList *SkipList) Remove(key int) {
	update := make(map[int]*Node)
	node := skipList.Header
	for i := skipList.Level - 1; i >= 0; i-- {
		for {
			if node.Forward[i].Value == nil {
				break
			}
			if node.Forward[i].Value.(int) == key {
				update[i] = node
				break
			}
			if node.Forward[i].Value.(int) < key {
				node = &node.Forward[i]
				continue
			} else {
				break
			}
		}
	}
	for i, v := range update {
		if v == skipList.Header {
			skipList.Level--
		}
		v.Forward[i] = v.Forward[i].Forward[i]
	}
}

func (skipList *SkipList) RandomLevel() int {
	level := 1
	for {
		if rand.Intn(2) == 1 { //产生0 1随机数
			level++
			if level >= SKIPLIST_MAX_LEVEL {
				break
			}
		} else {
			break
		}
	}

	return level
}

func (skipList *SkipList) PrintSkipList() {

	for i := SKIPLIST_MAX_LEVEL - 1; i >= 0; i-- {
		fmt.Println("level:", i)
		node := skipList.Header.Forward[i]
		for {
			if node.Value != nil {
				fmt.Printf("%d ", node.Value.(int))
				node = node.Forward[i]
			} else {
				break
			}
		}
		fmt.Println("\n-----------------------------")
	}

	fmt.Println("Current MaxLevel:", skipList.Level)
}
