# Fourinone
A distributed system, based on Fourinone 2.0
## Schedule
- Simply finish the distributed communication

## some ports
- 8000 : listen to register for workers
- 8001 : listen to hire workers 
- 8002 : server for cache
- 8003 : fs server

## How to run the demo?
### distribute calculate
- go run parkServerTest.go
- go run workTset.go IP PORT
- go run contractorTest.go

### distribute cache
- go run parkServerTest.go
- go run cacheNodeTest.go ip port maxbytes_for_cache

## Notice
供远程RPC调用的方法首字母需要大写， 大写表示允许外部访问
RPC方法最后需要对*reply赋值，而不是对reply赋值
方法声明的reply类型最好和传入的reply类型保持一致