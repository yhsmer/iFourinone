# iFourinone
A distributed system, based on Fourinone 2.0

## some ports
- 8000 : listen to register for workers
- 8001 : listen to hire workers 
- 8002 : server for cache
- 8003 : fs server

## How to run the demo?
### distribute calculate
- go run parkServerTest.go
- go run workTest.go IP PORT
> you can start one or more workTest.go
- go run contractorTest.go

### distribute cache
- go run parkServerTest.go
- go run cacheNodeTest.go ip port maxbytes_for_cache
> you can start one or more cacheNodeTest.go


### distribute fileSystem
- go run parkServerTest.go
- go run FSNodeTest.go localhost 19000 /Users/yexingming/code/go/iFourinone/123
- go run FSNodeTest.go localhost 19001 /Users/yexingming/code/go/iFourinone/456
> you can start one or more FSNodeTest.go

## Notice
- 供远程RPC调用的方法首字母需要大写，大写表示允许外部访问
- RPC方法若需要返回值最后需要对*reply赋值，而不是对reply赋值
- 方法声明的reply类型最好和传入的reply类型保持一致
