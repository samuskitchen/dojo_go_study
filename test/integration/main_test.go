package integration

import (
	data "dojo_go_study/config/database"
	testDB "dojo_go_study/config/database/test"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"

	"dojo_go_study/config/api"
	apiTest "dojo_go_study/config/api/test"
)

// Is a reference to the main Application type. This is used for its database
// connection that it harbours inside of the type as well as the route definitions
// that are defined on the embedded handler.
var database *data.Data
var server *api.Server

// TestMain calls testMain and passes the returned exit code to os.Exit(). The reason
// that TestMain is basically a wrapper around testMain is because os.Exit() does not
// respect deferred functions, so this configuration allows for a deferred function.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

// testMain returns an integer denoting an exit code to be returned and used in
// TestMain. The exit code 0 denotes success, all other codes denote failure (1
// and 2).
func testMain(m *testing.M) int {
	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	database = testDB.Open()
	defer database.DB.Close()

	port := os.Getenv("DAEMON_PORT")
	server = apiTest.NewServerTest(port, database)

	return m.Run()
}