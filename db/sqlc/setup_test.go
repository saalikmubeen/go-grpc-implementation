package generated_db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/saalikmubeen/backend-masterclass-go/utils"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := utils.LoadConfig("../..")

	if err != nil {
		log.Fatal("cannot load config:", err)

	}

	// Create a new database connection
	conn, err := sql.Open(config.DBDriver, config.DBURI)

	testDB = conn

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// Create a new testQueries object of type Queries for our tests to use
	testQueries = New(conn)

	os.Exit(m.Run())
}
