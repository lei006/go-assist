package main

import (
	"fmt"
	"time"

	"github.com/lei006/go-assist/pkg/daemon"
	"github.com/lei006/go-assist/pkg/zlog"
)

func main() {

	daemon.StartDaemon("testApp", "testApp", "", RunApp, zlog.New("daemon"))

}

func RunApp() {

	for true {
		fmt.Println("xxxxxxx")
		time.Sleep(time.Second)
	}

}
