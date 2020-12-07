package main

import (
	"github.com/seachenjy/go-comment/api"
	"github.com/seachenjy/go-comment/config"
)

func main() {
	err := config.Init("./config.yaml")
	if err != nil {
		panic(err)
	}

	api.Init()
}
