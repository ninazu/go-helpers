package main

import (
	"../../go-helpers"
	"fmt"
)

func main() {
	//TODO REMOVE KEY

	id, err := ninazu.ReCaptcha("123").
		GetCaptchaSolution("6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", "https://www.google.com/recaptcha/api2/demo")

	if err != nil {
		panic(err)
	}

	fmt.Println(id)
}
