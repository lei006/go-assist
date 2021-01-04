package main

import (
	"github.com/lei006/go-assist/zcli"

	"github.com/astaxie/beego"
)

func run() {

	beego.Run()
}

func stop() {

}

func debug_main() {
	run()
}

func server_main() {

	err := zcli.LaunchServiceRun("demo_test", "", run)
	if err == nil {
		stop()
	}
}

func main() {
	server_main()
}
