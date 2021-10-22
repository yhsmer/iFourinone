package main

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
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

func sayHelloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("qqqqq")
	http.Redirect(w, r, "http://localhost:8088", http.StatusFound) //重定向
	/*content :=[]byte("hello world")
	//向test.txt写入hello world
	err := ioutil.WriteFile("test.txt", content,0644)
	if err !=nil{
		panic(err)
	}*/
}

func ForwardHandler(writer http.ResponseWriter, request *http.Request) {
	u := &url.URL{
		Scheme: "https",
		Host:   "baidu.com",
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	request.URL.Path = ""
	proxy.ServeHTTP(writer, request)
}

func test() {
	p, _ := filepath.Abs(filepath.Dir("/Users/yexingming/Downloads/fourinone/"))
	http.Handle("/", http.FileServer(http.Dir(p)))
	http.HandleFunc("/a", ForwardHandler)//   设置访问路由
	log.Fatal(http.ListenAndServe(":8089",nil))
}

type student struct {
	Name string
	Age  int8
}

func ginDemo()  {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	//stu1 := &student{Name: "https://baidu.com", Age: 20}
	//stu2 := &student{Name: "Jack", Age: 22}
	str := []string{"http://localhost:9999/testfs", "https://baidu.com"}
	r.GET("/testfs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "testfs.html", gin.H{
			"title":  "Gin",
			"stuArr": str,
		})
	})

	r.Run(":9999")
}

func ginFs() {
	//ginDemo()
	//http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))

	r := gin.Default()
	r.LoadHTMLFiles( "./templates/123.html")
	// 静态文件服务
	// 显示当前文件夹下面的所有文件 / 或者指定文件
	// 页面返回：服务器当前路径下地文件信息
	r.StaticFS("/fs", http.Dir("/Users/yexingming/Downloads/fourinone/"))

	r.GET("/api/:name/:age", func(c *gin.Context) {
		name := c.Param("name")
		age := c.Param("age")
		c.JSON(http.StatusOK, gin.H{
			"name:": name,
			"age: ":age,
		})
	})


	r.Run("0.0.0.0:8000")

}

func uploadFile(c *gin.Context) {
	// FormFile方法会读取参数“123”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	file, header, err := c.Request.FormFile("123")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	// header调用Filename方法，就可以得到文件名
	filename := header.Filename
	fmt.Println(file, err, filename)

	// 获取创建文件的路径
	path := c.DefaultPostForm("path", "anonymous")
	fmt.Println(path)

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	out, err := os.Create( path + "/" + filename)

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	c.String(http.StatusCreated, "123 successful \n")
}

func main() {
	path := "./ioasdf/dsf/df"
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	fmt.Println(path)

	//router := gin.Default()
	//
	//// 调用POST方法，传入路由参数和路由函数
	//router.POST("/123", uploadFile)
	//
	//// 监听端口8000，注意冒号。
	//router.Run(":8000")
}

/*
curl -X POST http://127.0.0.1:8000/upload -F "upload=@/Users/yexingming/Pictures/txt" -H "Content-Type: multipart/form-data"
curl -X POST http://127.0.0.1:8000/upload -F "upload=@/Users/yexingming/Pictures/txt" -F "path=." -H "Content-Type: multipart/form-data"
curl -X POST http://127.0.0.1:8000/upload -F "upload=@/Users/yexingming/Pictures/txt" -F "path=./123" -H "Content-Type: multipart/form-data"
*/
