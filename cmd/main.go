package main

import (
	"dojo_go_study/config/api"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
)

func main() {
	log.Println("stating API cmd")
	port := os.Getenv("DAEMON_PORT")
	api.Start(port)
}
