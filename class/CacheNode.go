package class

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"
)
var wg = sync.WaitGroup{}

// RPC service
type CacheNode struct {
	addr string
	LruCache *LRUCache
}

func NewCacheNode(addr string, maxBytes int64) *CacheNode {
	c := CacheNode{}
	c.addr = addr
	c.LruCache = NewLRUCache(maxBytes)
	return &c
}

func (c CacheNode) getAddr() string {
	return c.addr	
}

func (c CacheNode) Get(key string) (Value,bool){
	value, ok := c.LruCache.Get(key)
	if ok{
		return value, ok
	}else{
		return "No such value",false
	}
}

func (c CacheNode) Add(key string, value Value) {
	c.LruCache.Add(key, value)
}

func startCacheRPC(ip string, port string, maxBytes int64) {
	//注册RPC服务
	cacheNode := NewCacheNode(ip + ":" + port, maxBytes)

	rpc.Register(cacheNode)

	//创建监听端口
	tcpAddr,err:=net.ResolveTCPAddr("tcp",ip+":"+port)
	if err!=nil{
		log.Println("startCacheRPC net.ResolveTCPAddr err:",err)
		return
	}
	listener,err:=net.ListenTCP("tcp",tcpAddr)
	if err!=nil{
		log.Println("startCacheRPC net.ListenTCP err:",err)
		return
	}
	defer listener.Close()

	//循环监听服务
	for{
		conn,err:=listener.Accept()
		if err!=nil{
			log.Println("startCacheRPC listener.Accept err:",err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func (c *CacheNode) GetCache(key string, value *Value) error {
	v, ok := c.Get(key)
	if ok{
		*value = v
	}
	return nil
}

func (c *CacheNode) AddCache(kv string, ret *Value) error{
	strs := strings.Split(kv, " ")
	c.Add(strs[0], strs[1])
	return nil
}

func (c CacheNode) CacheStart() {
	//获取命令行参数
	args := os.Args
	if len(args) != 4 {
		log.Println("Please input your hostIP, running Port and maxBytes for cacheNode!")
		log.Println("The right form is ：go run xxx.go IP Port maxBytes_for_cache")
		return
	}
	ip := args[1]
	port := args[2]
	maxBytes,_ := strconv.ParseInt(args[3], 10, 64)

	//开启RPC服务
	go startCacheRPC(ip, port, maxBytes)

	//向职介者注册
	go logInToPark(ip, port, "cache")

	for{}
}



