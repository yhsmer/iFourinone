package consistenthash

import (
	"fmt"
	"strconv"
"testing"
)

func TestHashing(t *testing.T) {
	hash := New(func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"4": "2",
		"5": "4",
		"7": "2",
	}

	for k, v := range testCases {
		fmt.Println(k, hash.Get(k))
		if hash.Get(k) != v {

		}
	}

	// Adds 8, 18, 28
	hash.Add("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		fmt.Println(k,hash.Get(k))
		if hash.Get(k) != v {

		}
	}

}