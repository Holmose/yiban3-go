package main

import (
	"Yiban3/clockrpc"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	err := rpc.Register(clockrpc.ClockService{})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("ClockOne Server Listening 0.0.0.0:29760 ...")
	listen, err := net.Listen("tcp", ":29760")
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// 正在执行时不接收其他请求
		jsonrpc.ServeConn(conn)
	}
}
