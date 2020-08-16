package db_storage

import (
	//"database/sql"
	//"github.com/sirupsen/logrus"
	//"strings"

	"database/sql/driver"
	//"fmt"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)


type sport_pair struct {
	sportName string
	sportRatio string
}

var tests = []sport_pair {
	{ "soccer", "3.123", },
	{ "football", "2.549", },
	{ "baseball", "1.091", },
	{ "soccer", "4.129", },
	{ "baseball", "1.901", },
}

func TestPutSportLine(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	for _, pair := range tests {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO public\."` + pair.sportName +`"`).WithArgs(
			driver.Value(pair.sportName),
			driver.Value(pair.sportRatio),
		).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	for _, pair := range tests {
		PutSportLine(db, pair.sportName, pair.sportRatio)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expections: %s", err)
	}
}

func TestGetSportRatio(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expectedRows := sqlmock.NewRows([]string{"sportratio"}).AddRow("3.123")

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT "sportratio" FROM`).WillReturnRows(expectedRows)
	mock.ExpectCommit()

	GetSportRatio(db, "soccer")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expections: %s", err)
	}
}