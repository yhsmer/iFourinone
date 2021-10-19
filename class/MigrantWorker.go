package class

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
)

//农民工
type Workers struct {
}

// RPC服务
// 本地用于给远程调用的方法需要注册成一个RPC服务之后才能给远程使用
type WorkRPC struct {
	ip, port string
	ready bool
}

//向职介者注册
func logInToPark(ip string, port string, logType string){
	//创建TCP连接,连接职介者在该端口的农名工监听服务
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Println("logInToPark net.Dial() err:", err)
		return
	}
	defer conn.Close()

	//发送注册的RPC服务地址
	// 农名工工作在哪个IP哪个端口
	conn.Write([]byte(ip + ":" + port + "-" + logType))

	//接收来自ParkServer的注册成功信息
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("logInToPark conn.Read err:", err)
		return
	}
	if string(buf[:n]) == "ok" {
		log.Println("log in successfully , waiting work .....")
	}
}

//开启RPC服务
func startRPC(ip string, port string) {
	//注册RPC服务
	//workRPC := &WorkRPC{ip,port,true}
	workRPC := new(WorkRPC)
	workRPC.ip = ip
	workRPC.port = port
	workRPC.ready = true

	rpc.Register(workRPC)

	//创建监听端口
	tcpAddr,err:=net.ResolveTCPAddr("tcp",ip+":"+port)
	if err!=nil{
		log.Println("startRPC net.ResolveTCPAddr err:",err)
		return
	}
	listener,err:=net.ListenTCP("tcp",tcpAddr)
	if err!=nil{
		log.Println("startRPC net.ListenTCP err:",err)
		return
	}
	defer listener.Close()

	//循环监听服务
	for{
		conn,err:=listener.Accept()
		if err!=nil{
			log.Println("startRPC listener.Accept err:",err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// 供远程调用的让节点工作的方法
func (workRPC *WorkRPC) DoTask(args *string, ret *map[string]int) error {
	log.Println(workRPC.ip + workRPC.port + ": is working")

	strs := strings.Fields(strings.TrimSpace(*args))
	for _, s := range strs{
		(*ret)[s]++
		log.Println(s)
	}

	workRPC.ready = true

	return nil
}

// 供远程调用查看节点状态
func (workRPC *WorkRPC) CheckStatus(redundant int, ret *bool) error {
	log.Println(workRPC.ready)
	*ret = workRPC.ready
	return nil
}


//农民工启动
func (workers Workers) StartWork() {
	//获取命令行参数
	args := os.Args
	if len(args) != 3 {
		log.Println("Please input your hostIP and running Port!")
		log.Println("The right form is ：go run xxx.go IP Port")
		return
	}
	ip := args[1]
	port := args[2]

	//开启RPC服务
	go startRPC(ip,port)

	//向职介者注册
	go logInToPark(ip, port, "worker")

	for {}
}
