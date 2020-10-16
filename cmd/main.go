package main

import (
	"dojo_go_study/config/api"
	"dojo_go_study/config/database"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("error in main")

			os.Exit(1)
		}
	}()

	// connection to the database.
	db := database.New()
	if err := db.DB.Ping(); err != nil {
		log.Fatal(err)
	}

	conn := &database.Data{
		DB: db.DB,
	}

	DaemonPort := os.Getenv("DAEMON_PORT")
	serv := api.NewApplication(DaemonPort, conn)

	// start the server.
	go serv.Start()

	// Wait for an in interrupt.
	// If you ask about "<-" look here https://tour.golang.org/concurrency/2
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown.
	_ = serv.Close()
	_ = database.Close()
}