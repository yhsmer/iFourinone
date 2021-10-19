package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"unsafe"
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


type Contractor struct {
	name string
	ctor *Contractor
}

func (contractor Contractor)toNext(ctor Contractor) *Contractor{
	fmt.Println(contractor.name)
	contractor.ctor = &ctor
	return &ctor
}

func (contractor Contractor) printName(){
	fmt.Println(contractor.name)
}

type Cache struct {
	MaxBytes   int64 // 允许使用的最大内存
	NowBytes   int64 //当前已经使用的内存
	data       *list.List
	cache      map[string]*list.Element
	mu         sync.Mutex
}

type date struct {
	key   string
	value Value
}
type Value interface {
}

func NewCache(maxBytes int64) *Cache {
	return &Cache{
		MaxBytes:  maxBytes,
		NowBytes:  0,
		data:      list.New(),
		cache:     make(map[string]*list.Element),
	}
}

// Get get the key's value, move the node to the front of the queue, means its frequently visited
func (c *Cache) Get(key string) (value Value, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, ok := c.cache[key]; ok {
		c.data.MoveToFront(ele)
		kv := ele.Value.(*date)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) removeOldest() {
	ele := c.data.Back()
	if ele != nil {
		c.data.Remove(ele)
		kv := ele.Value.(*date)
		delete(c.cache, kv.key)
		c.NowBytes -= int64(unsafe.Sizeof(kv.key)) + int64(unsafe.Sizeof(kv.value))
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if date exist, then replace
	if ele, ok := c.cache[key]; ok {
		c.data.MoveToFront(ele)
		kv := ele.Value.(*date)
		c.NowBytes += int64(unsafe.Sizeof(value)) - int64(unsafe.Sizeof(kv.value))
		kv.value = value
	} else {
		ele := c.data.PushFront(&date{key, value})
		c.cache[key] = ele
		c.NowBytes += int64(unsafe.Sizeof(key)) + int64(unsafe.Sizeof(value))
	}
	for c.MaxBytes != 0 && c.MaxBytes < c.NowBytes {
		c.removeOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.data.Len()
}
type server int

func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/cache") {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
		//panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	if strings.Contains(r.URL.Path, "Get"){
		log.Println("Get....")
	}else if strings.Contains(r.URL.Path, "Add"){
		log.Println("Add...")
	}



	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	log.Println(key)
	log.Println(value)
	log.Println(r.URL.Path)
	w.Write([]byte("Hello World!"))
}

func startServer() {
	var s server
	http.ListenAndServe("localhost:9999", &s)
}



func main() {
	//startServer()
	var v Value
	v = 123
	k := fmt.Sprint(v)
	fmt.Println(k)
}
