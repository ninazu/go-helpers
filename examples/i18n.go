package main

import (
	"../../go-helpers/i18n"
	"github.com/ninazu/go-helpers"
	"fmt"
)

func main() {
	i18n.Init(ninazu.PATH_OF_SOURCE()+"/../messages", "ru")
	
	err := i18n.Generate(ninazu.PATH_OF_SOURCE(), []string{"en", "fr"})
	
	if err != nil {
		panic(err)
	}
	
	fmt.Println(i18n.Translate("Hello"))
}
