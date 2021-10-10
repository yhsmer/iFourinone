package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type WareHouse interface{
	
}


func ReadFile(file_name string) (info string) {
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var lineText string
	var ans string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText = scanner.Text()
		log.Println(">>>>>" + lineText)
		ans += lineText + " "
	}

	return ans
}

func fun(m map[string]string)int{
	file_name := m["filepath"]

	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ret := []string{"1","2","3"}

	worksInput := make([]string, len(ret))
	i := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		worksInput[i] += scanner.Text() + " "
		i = (i+1)%len(ret)
	}

	for _, v := range worksInput{
		log.Println(v)
		log.Println(">>>")
	}

	worksOutput := make([]map[string]int, len(ret))
	for i := 0; i < len(worksOutput); i++{
		worksOutput[i] = make(map[string]int)
	}

	DoTask(&worksInput[0], &worksOutput[0])

	for k,v := range worksOutput[0]{
		log.Println(k + ":" + strconv.Itoa(v))
	}

	return 0
}

func DoTask(args *string, ret *map[string]int) int {
	strs := strings.Fields(strings.TrimSpace(*args))
	for _, s := range strs{
		(*ret)[s]++
		log.Println(s)
	}
	return 1
}


func main(){
	m := map[string]string{"filepath": "./in"}
	fun(m)
}
