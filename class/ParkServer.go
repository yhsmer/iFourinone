package class

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"iFourinone/consistenthash"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"
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
	caches []string
	fs []string
	consistendHash *consistenthash.Consistenthash
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
	log.Println("ParkServer is waiting workers's connection at localhost:8000...")

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
	args := strings.Split(string(buf[:n]),"-")

	if args[1] == "worker" {
		works.WorkAddr = args[0]
		park.workAddrs = append(park.workAddrs, works.WorkAddr)

		/*log.Println(works.WorkAddr, "connect successfully!")
		log.Println("success to get the", works.WorkAddr, "'s log message: I'm", string(buf[:n]))*/

	} else if args[1] == "cache"{
		if park.consistendHash == nil{
			park.consistendHash = consistenthash.New(nil)
		}
		park.consistendHash.Add(args[0])
		park.caches = append(park.caches, args[0])

	} else if args[1] == "fs"{
		park.fs = append(park.fs, "http://" + args[0])
	}

	log.Println(args[0], "connect successfully!")
	log.Println("success to get the", args[1], "'s log message: I'm", args[0])

	//回写表示接收成功
	conn.Write([]byte("ok"))
}

type server struct{
	// http://localhost:8002/${bashPath}
	basePath string
}

func NewServer() *server {
	return &server{
		basePath: "/cache",
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, s.basePath) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	log.Println(r.URL.Path)

	if strings.Contains(r.URL.Path, "Get"){
		key := r.URL.Query().Get("key")
		log.Println("Get.... key -> " + key)

		addr := park.consistendHash.Get(key)
		if addr == ""{
			w.Write([]byte("There has no cache node existed"))
			return
		}

		// 创建TCP连接,连接职介者在该端口的监听工头雇工服务
		conn, err := rpc.Dial("tcp", addr)
		if err != nil {
			log.Println("ServerHttp Dial rpc.Dial err:", err)
			return
		}
		// 在当前函数返回前执行传入的函数，一般用于关闭资源等等
		defer conn.Close()

		//远程调用ParkServer的方法
		var ret Value
		err = conn.Call("CacheNode.GetCache", key, &ret)
		// ret返回之后为string
		if err != nil {
			log.Println("CacheNode.GetCache conn.Call err:", err)
			return
		}
		// transform Value to string
		value := fmt.Sprint(ret)
		log.Println("key -> " + key + " & value -> " + value)
		w.Write([]byte(value))

	} else if strings.Contains(r.URL.Path, "Add") {
		log.Println("Add....")
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		log.Println("Add : key -> " + key + " value -> " + value)

		addr := park.consistendHash.Get(key)
		if addr == ""{
			w.Write([]byte("There has no cache node existed"))
			return
		}

		// 创建TCP连接,连接职介者在该端口的监听工头雇工服务
		conn, err := rpc.Dial("tcp", addr)
		if err != nil {
			log.Println("ServerHttp Dial rpc.Dial err:", err)
			return
		}
		// 在当前函数返回前执行传入的函数，一般用于关闭资源等等
		defer conn.Close()

		//远程调用ParkServer的方法
		var ret Value
		err = conn.Call("CacheNode.AddCache", key + " " + value, &ret)
		// ret返回之后为string
		if err != nil {
			log.Println("CacheNode.AddCache conn.Call err:", err)
			return
		}
		w.Write([]byte("Add cache successful"))
	}
}

func startFsServer()  {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/fs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "fs.html", gin.H{
			"title":  "All Servers:",
			"stuArr": park.fs,
		})
	})
	log.Println("fs server is running at localhost:8003")

	err := r.Run(":8003")
	if err != nil {
		log.Println("startFsServer r.Run error")
		log.Println(err)
		return
	}
}

//启动ParkServer
func (park Park) ParkStart() {
	// 启动监听农民工注册服务
	go logInServer()

	// 启动监听包工头雇工服务
	go waitingWorkerServer()

	go startFsServer()

	s := NewServer()
	log.Println("cache server is running at localhost:8002")
	err := http.ListenAndServe("localhost:8002", s)
	if err != nil {
		log.Println("http.ListenAndServe (port : 8002) err:", err)
		return
	}
	for{}
}
