package main

import (
	"fmt"

	"github.com/lei006/go-assist/license"
)

func main() {
	key := "341234123"

	data := license.LicenseData{}
	data.StandardClaims.Id = "aaa"

	sin, err := data.Sign(key)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	data11 := license.LicenseData{}

	data11.Verify(sin, key)

	fmt.Println(data11)
}
