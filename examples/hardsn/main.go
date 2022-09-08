package main

import (
	"fmt"

	"github.com/lei006/go-assist/hardsn"
)

func main() {
	id_md5, err := hardsn.GetPhysicalID()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("hardsn:", id_md5)
	//test_license(100)

}
