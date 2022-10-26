package main

import (
	"Yiban3/clockrpc"
	"fmt"
	"net"
	"net/rpc/jsonrpc"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:29760")
	if err != nil {
		panic(err)
	}
	client := jsonrpc.NewClient(conn)
	var result string
	err = client.Call("ClockService.Clock",
		clockrpc.Args{Key: "hFjCM5XBMC6bo3k", Iv: "hONmvJHk"}, &result)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(result)
	}
}
