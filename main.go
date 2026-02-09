package main

import (
	"github.com/mundotv789123/raspadmin/cmd"
	"github.com/mundotv789123/raspadmin/config"
)

func main() {
	config.Init()
	cmd.Run()
}
