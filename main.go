package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mundotv789123/raspadmin/cmd"
)

func main() {
	envFile := ".env"
	if err := godotenv.Load(envFile); err != nil {
		log.Print("Error loading .env file: ", err)
	}
	cmd.Run()
}
