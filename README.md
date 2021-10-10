# Fourinone
A distributed system, based on Fourinone 2.0
## Schedule
- Simply finish the distributed communication
## How to run the demo?
- go run parkServerTest.go
- go run workTset.go IP PORT
- go run contractorTest.go
## Notice
供远程RPC调用的方法首字母需要大写， 大写表示允许外部访问
RPC方法最后需要对*reply赋值，而不是对reply赋值
方法声明的reply类型最好和传入的reply类型保持一致