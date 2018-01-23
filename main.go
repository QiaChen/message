package main

func main() {
	
	initSys()
	//flash跨域接口
	go CrossDomainServer()
	//启动接收信号进程
	receive()
}

func initSys(){
	Topics = make(map[string]*topic)
	UsersTopocs = make(map[string]map[string]string)
}
