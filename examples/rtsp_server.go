package main

import (
	"fmt"

	"go-assist/rtsp_server"
)

func main() {
	fmt.Println("1111111111")
	go rtsp_server.GetServer().Start()

}
