package api

import (
	"dojo_go_study/config/database"
)

func Start(port string)  {

	// connection to the database.
	db := database.New()
	defer db.DB.Close()

	server := newServer(port, db)

	// start the server.
	server.Start()
}
