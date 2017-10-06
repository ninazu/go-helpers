package main

import (
	"github.com/ninazu/go-helpers"
	"fmt"
)

func main() {
	r := ninazu.ReCaptcha("1abc234de56fab7c89012d34e56fa7b8").
		GetCaptcha("6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-", "https://www.google.com/recaptcha/api2/demo")

	fmt.Println(r)
}
