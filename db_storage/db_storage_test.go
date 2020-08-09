package db_storage

import (
	_ "github.com/lib/pq"
	"testing"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwerty"
	dbname   = "LinesStorage"
)

func TestStorageInit(t *testing.T)  {
	StorageInit(host, port, user, password, dbname)
}

func TestIsConnected(t *testing.T) {
	db := StorageInit(host, port, user, password, dbname)
	IsConnected(db)
}


type testSportRatio struct {
	sportName string
	ratioValue string
}

var testValues = []testSportRatio {
	{ "SOCCER", "1.325"},
	{ "FOOTBALL", "0.321"},
	{ "BASEBALL", "3.213"},
	{ "SOCCER", "2.123"},
	{ "SOCCER", "0.369"},
	{ "BASEBALL", "1.461"},
	{ "SOCCER", "0.549"},
	{ "BASEBALL", "2.1227"},
	{ "SOCCER", "0.005"},
	{ "FOOTBALL", "0.1"},
	{ "FOOTBALL", "2.001"},
}

func TestGetPutSportline(t *testing.T) {
	db := StorageInit(host, port, user, password, dbname)
	for _, pair := range testValues {
		PutSportLine(db, pair.sportName, pair.ratioValue)
		testValue := GetSportRatio(db, pair.sportName)
		if testValue != pair.ratioValue {
			t.Error(
				"Pushed to DB: ", pair.sportName, pair.ratioValue,
				"Pulled after:", pair.sportName, testValue,
			)
		}
	}
}