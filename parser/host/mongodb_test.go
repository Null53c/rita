package host

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/globalsign/mgo/dbtest"
	"github.com/activecm/rita/database"
	"github.com/activecm/rita/parser/parsetypes"
)

// Server holds the dbtest DBServer
var Server dbtest.DBServer

// Set the test database
var testTargetDB = "tmp_test_db"

var testRepo Repository

func TestCreateIndexes(t *testing.T) {
	err := testRepo.CreateIndexes(testTargetDB)
	if err != nil {
		t.Errorf("Error creating host indexes")
	}
}

func TestUpsert(t *testing.T) {
	testHost := &parsetypes.Host{
		IP: "127.0.0.1",
		Local: true,
		IPv4: true,
		CountSrc: 123,
		CountDst: 123,
		IPv4Binary: 123,
		MaxDuration: 123.0,
		MaxBeaconScore: 123.0,
		MaxBeaconConnCount: 123,
		BlOutCount: 123,
		BlInCount: 123,
		BlSumAvgBytes: 123,
		BlTotalBytes: 123,
		TxtQueryCount: 123,
	}

	err := testRepo.Upsert(testHost, true, testTargetDB)
	if err != nil {
		t.Errorf("Error upserting host")
	}
}

// TestMain wraps all tests with the needed initialized mock DB and fixtures
func TestMain(m *testing.M) {
	// Store temporary databases files in a temporary directory
	tempDir, _ := ioutil.TempDir("", "testing")
	Server.SetPath(tempDir)

	// Set the main session variable to the temporary MongoDB instance
	ssn := Server.Session()

	db := database.DB{Session: ssn}

	testRepo = NewMongoRepository(&db)

	// Run the test suite
	retCode := m.Run()

	// Clean up test database and session
	ssn.DB(testTargetDB).DropDatabase()
	ssn.Close()

	// Shut down the temporary server and removes data on disk.
	Server.Stop()

	// call with result of m.Run()
	os.Exit(retCode)
}