package main

import (
	"fmt"
	"github.com/ninazu/go-helpers"
)

func main() {
	token, err := ninazu.ReCaptcha("YOUR_KEY").
		GetCaptchaSolution("6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", "https://www.google.com/recaptcha/api2/demo")

	if err != nil {
		panic(err)
	}

	fmt.Println(token)
}
