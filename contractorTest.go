package main

import (
	"bufio"
	"fmt"
	"iFourinone/class"
	"log"
	"os"
	"strconv"
)

type ContractorTest struct {
	class.Contractor
}

func (c ContractorTest) ScheduleWork(m map[string]string) {
	ret := c.Contractor.GetAllWorkers()

	file_name := m["filepath"]

	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	worksInput := make([]string, len(ret))
	i := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		worksInput[i] += scanner.Text() + " "
		i = (i+1)%len(ret)
	}

	worksOutput := c.Contractor.AssignWork(worksInput,ret)

	finalOutput := make(map[string]int)
	for _, m := range worksOutput{
		for k,v := range m{
			finalOutput[k] += v
		}
	}

	fmt.Println("Contractor's output is as follows")
	for k,v := range finalOutput{
		fmt.Println(k + " " + strconv.Itoa(v))
	}
	return
}


func main() {
	c := ContractorTest{}
	//var c ContractorTest
	//var contractor class.Contractor
	m := map[string]string{"filepath": "./in"}

	c.ScheduleWork(m)

	// 工头之间的串行处理
	//contractor.ToNext(class.Contractor{Name: "TwoCtor"},m).ToNext(class.Contractor{Name: "ThreeCtor"},m)

}
