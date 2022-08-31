package main

import (
	"fmt"

	"github.com/powenyu/split-order/bot"
	"github.com/powenyu/split-order/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
