package main

import (
	"github.com/mundotv789123/raspadmin/cmd"
	"github.com/mundotv789123/raspadmin/internal/config"
)

func main() {
	config.Init()
	cmd.Run()
}
