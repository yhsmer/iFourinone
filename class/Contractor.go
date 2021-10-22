package class

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

type Contractor struct {
	Name string
	Ctor *Contractor
}

func (contractor Contractor) ToNext(ctor Contractor, m map[string]string) *Contractor{
	//fmt.Println(contractor.name)
	contractor.Ctor = &ctor
	ctor.ScheduleWork(m)
	return &ctor
}

func (contractor Contractor) GetAllWorkers() []string {
	// 创建TCP连接,连接职介者在该端口的监听工头雇工服务
	conn, err := rpc.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		log.Println("QueryAllWorkers rpc.Dial err:", err)
		return nil
	}
	// 在当前函数返回前执行传入的函数，一般用于关闭资源等等
	defer conn.Close()

	//远程调用ParkServer的方法
	var ret []string
	conn.Call("Service.QueryAllWorkers", 1, &ret)
	// ret返回之后为string
	if err != nil {
		log.Println("QueryAllWorkers conn.Call err:", err)
		return nil
	}
	log.Printf("All waiting works: %v", ret)
	return ret
}

func (contractor Contractor) ScheduleWork(m map[string]string) {
	/*// 创建TCP连接,连接职介者在该端口的监听工头雇工服务
	conn, err := rpc.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		log.Println("QueryAllWorkers rpc.Dial err:", err)
		return
	}
	// 在当前函数返回前执行传入的函数，一般用于关闭资源等等
	defer conn.Close()

	//远程调用ParkServer的方法
	var ret []string
	conn.Call("Service.QueryAllWorkers", 1, &ret)
	// ret返回之后为string
	if err != nil {
		log.Println("QueryAllWorkers conn.Call err:", err)
		return
	}
	log.Printf("All waiting works: %v", ret)*/

	ret := contractor.GetAllWorkers()

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

	worksOutput := contractor.AssignWork(worksInput,ret)

	/*
	// system call
	worksOutput := make([]map[string]int, len(ret))
	for i := 0; i < len(worksOutput); i++{
		worksOutput[i] = make(map[string]int)
	}

	calls := make([]*rpc.Call,len(ret))
	for i, s := range ret{
		client, err := rpc.Dial("tcp", s)
		if err != nil {
			log.Fatal("DialWorker err :", err)
			continue
		}

		//var reply bool
		calls[i] = client.Go("WorkRPC.DoTask", &worksInput[i], &worksOutput[i],nil)
		if calls[i].Error != nil {
			log.Fatal("DoTask Call err", err)
		} else{
			log.Println("worker " + s + " assigned successfully")
		}
	}

	for {
		flag := true
		for _, c := range calls {
			if k := <-c.Done; k.Error != nil {
				log.Println("There is still exist one or more worker don't have finished the work")
				flag = false
			}
		}

		if flag == true{
			log.Println("Yes! worker ares done")
			break
		}
	}*/

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
}

func (contractor Contractor) AssignWork(worksInput []string, ret []string) []map[string]int{
	// system call
	worksOutput := make([]map[string]int, len(ret))
	for i := 0; i < len(worksOutput); i++{
		worksOutput[i] = make(map[string]int)
	}

	calls := make([]*rpc.Call,len(ret))
	for i, s := range ret{
		client, err := rpc.Dial("tcp", s)
		if err != nil {
			log.Fatal("DialWorker err :", err)
			continue
		}

		//var reply bool
		calls[i] = client.Go("WorkRPC.DoTask", &worksInput[i], &worksOutput[i],nil)
		if calls[i].Error != nil {
			log.Fatal("DoTask Call err", err)
		} else{
			log.Println("worker " + s + " assigned successfully")
		}
	}

	for {
		flag := true
		for _, c := range calls {
			if k := <-c.Done; k.Error != nil {
				log.Println("There is still exist one or more worker don't have finished the work")
				flag = false
			}
		}

		if flag == true{
			log.Println("Yes! worker ares done")
			break
		}
	}
	return worksOutput
}