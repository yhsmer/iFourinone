package main

import (
	"iFourinone/class"
	"log"
	"strings"
)

func main() {
	var works class.Workers
	d := class.New(func(args *string, ret *map[string]int) *map[string]int {
		strs := strings.Fields(strings.TrimSpace(*args))
		for _, s := range strs{
			(*ret)[s]++
			log.Println(s)
		}
		return ret
	})

	works.StartWork(d)
}
