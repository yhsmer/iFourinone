package class

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestCache_Get(t *testing.T) {
	node := NewLRUCache(10240)
	node.Add("1","yexm")
	fmt.Println(node.Get("1"))
}

func TestCache_Add(t *testing.T) {
	node := NewLRUCache(10240)
	node.Add("1","yexm")
	fmt.Println(node.NowBytes)
	fmt.Println(unsafe.Sizeof("1"))
	fmt.Println(unsafe.Sizeof("yexm"))
	if node.NowBytes != int64(unsafe.Sizeof("1") + unsafe.Sizeof("yexm")){
		t.Fatal("error")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	node := NewLRUCache(10240)
	node.Add("1","yexm")
	node.Add("2",123)
	node.removeOldest()
	fmt.Println(node.data.Front().Value)
}
