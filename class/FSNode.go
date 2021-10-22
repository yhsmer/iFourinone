package class

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type FSNode struct {

}

var path string

func uploadFile(c *gin.Context) {
	// FormFile方法会读取参数“123”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	// header调用Filename方法，就可以得到文件名
	filename := header.Filename
	fmt.Println(file, err, filename)

	// 获取创建文件的路径
	cpath := c.DefaultPostForm("path", ".")
	if strings.HasSuffix(cpath, "/") {
		cpath = cpath[:len(cpath)-1]
	}
	fmt.Println(cpath)

	// path = /usr/test
	// cpath = .
	if cpath == "."{
		cpath = path
	} else {
		// cpath = ./go
		cpath = path + "/" + cpath
	}

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	out, err := os.Create( cpath + "/" + filename)

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	c.String(http.StatusCreated, "upload successful \n")
}

func deleteFile(c *gin.Context) {
	/*
		curl -X POST http://127.0.0.1:19000/delete -F "file=1.txt" -H "Content-Type: multipart/form-data"
		curl -X POST http://127.0.0.1:19000/delete -F "file=dir/1.txt" -H "Content-Type: multipart/form-data"
	*/

	// 获取需要删除的文件路径
	file := c.DefaultPostForm("file", "")
	if strings.HasSuffix(file, "/") {
		file = file[:len(file)-1]
	}

	if file == ""{
		return
	}
	fmt.Println(file)

	// path = /usr/test
	// file = 1.txt
	// file = dir/1.txt
	file = path + "/" + file


	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	err := os.Remove(file)

	if err != nil {
		log.Fatal(err)
	}

	c.String(http.StatusCreated, "delete successful \n")
}



func startServer(ip, port, path string)  {
	router := gin.Default()
	// 静态文件服务
	// 显示当前文件夹下面的所有文件 / 或者指定文件
	// 页面返回：服务器当前路径下地文件信息
	router.StaticFS("", http.Dir(path))

	// 上传文件服务
	router.POST("/upload", uploadFile)
/*
	curl -X POST http://127.0.0.1:19000/upload -F "upload=@/Users/yexingming/Pictures/1.txt" -F "path=dir" -H "Content-Type: multipart/form-data"
	curl -X POST http://127.0.0.1:19000/upload -F "upload=@/Users/yexingming/Pictures/1.txt" -F "path=." -H "Content-Type: multipart/form-data"
*/
	router.POST("/delete", deleteFile)
/*
	curl -X POST http://127.0.0.1:19000/delete -F "file=1.txt" -H "Content-Type: multipart/form-data"
	curl -X POST http://127.0.0.1:19000/upload -F "file=dir/1.txt" -H "Content-Type: multipart/form-data"
*/

	router.Run(":" + port)
}

func (fs FSNode)FSStart() {
	//获取命令行参数
	args := os.Args
	if len(args) != 4 {
		log.Println("Please input your hostIP, running Port and dir_path for FSNode!")
		log.Println("The right form is ：go run xxx.go IP Port dir_path")
		return
	}
	
	ip := args[1]
	port := args[2]
	path = args[3]

	// 去除path末尾的/
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}


	go startServer(ip, port, path)

	go logInToPark(ip, port, "fs")

	for{}
}

/*
   upload commands:
   	curl -X POST http://127.0.0.1:19005/upload -F "upload=@/Users/yexingming/Pictures/1.txt" -F "path=dir" -H "Content-Type: multipart/form-data"
   	curl -X POST http://127.0.0.1:19005/upload -F "upload=@/Users/yexingming/Pictures/1.txt" -F "path=." -H "Content-Type: multipart/form-data"

   delete commands:
	curl -X POST http://127.0.0.1:19005/delete -F "file=dir/1.txt" -H "Content-Type: multipart/form-data"
   	curl -X POST http://127.0.0.1:19005/delete -F "file=1.txt" -H "Content-Type: multipart/form-data"
*/




