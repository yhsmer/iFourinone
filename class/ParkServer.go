package class

import (
	"log"
	"net"
	"net/rpc"
)

//农民工注册信息结构
type WaitingWorkers struct {
	//workerIP   string
	//workerPort string
	//Ready      bool
	WorkAddr string
}

//职介者信息,用于启动Park进程，启动之后直接调用函数，而不是方法
type Park struct {
	// 存储注册信息汇总
	workAddrs []string
}
var park Park

//RPC调用 远程过程调用
type Service struct{}

//RPC for Contractor
func (s *Service) QueryAllWorkers(redundant int, ret *[]string) error {
	var readyNodes []string
	for _, workAddr := range park.workAddrs{
		log.Println("All waiting works: " + workAddr)
		client, err := rpc.Dial("tcp", workAddr)
		if err != nil {
			log.Fatal("QueryAllWorkers dialing err :", err)
			continue
		}

		var reply bool
		err = client.Call("WorkRPC.CheckStatus", 1, &reply)
		if err != nil {
			log.Fatal("QueryAllWorkers Call err", err)
			continue
		}
		if reply == true{
			readyNodes = append(readyNodes, workAddr)
		}
	}
	*ret = readyNodes
	return nil
}
/*
client, err := rpc.Dial("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }

    var reply string
    err = client.Call("HelloService.Hello", "hello", &reply)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(reply)
 */


//循环监听包工头雇工服务，设置监听，提供RPC服务的函数在此注册
//正常连接之后通过rpc.ServeConn函数在该TCP上使用注册的（提供RPC服务）函数为对方提供服务
func waitingWorkerServer() {
	//注册RPC服务
	s := new(Service)
	rpc.Register(s)

	//创建监听端口
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8001")
	if err != nil {
		log.Println("waitingWorkerServer net.ResolveTCPAddr err:", err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println("waitingWorkerServer net.Listen err:", err)
		return
	}
	defer listener.Close()

	//循环接收服务
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("waitingWorkerServer listener.Accept err:", err)
			continue
		}
		// 通过rpc.ServeConn函数在该TCP链接上为对方提供RPC服务。
		go rpc.ServeConn(conn)
	}
}

//循环监听MigrantWorker的注册
func logInServer() {
	//创建监听socket;固定IP+端口为"127.0.0.1:8000"
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Println("logInServer net.Listen err:", err)
		return
	}
	defer listener.Close()
	log.Println("ParkServer is waiting connection ...")

	//循环监听注册
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("logInServer listener.Accept err:", err)
			return
		}
		go HandleLogInConn(conn)
	}
}

//用于处理注册
func HandleLogInConn(conn net.Conn) {
	defer conn.Close()

	//ParkServer接收注册信息
	var works WaitingWorkers
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("HandleConn conn.Read err:", err)
		return
	}

	works.WorkAddr = string(buf[:n])
	park.workAddrs = append(park.workAddrs, works.WorkAddr)

	log.Println("MigrantWorker of", works.WorkAddr, "connect successfully!")

	log.Println("success to get the", works.WorkAddr, "'s log message: I'm", string(buf[:n]))

	//回写表示接收成功
	conn.Write([]byte("ok"))

}

//启动ParkServer
func (park Park) ParkStart() {

	//启动监听农民工注册服务
	go logInServer()

	//启动监听包工头雇工服务
	go waitingWorkerServer()

	for{}
}
