package http_workers

import (
	"LineProcessor/db_storage"
	"database/sql"
	"database/sql/driver"
	"strings"
	"testing"
	"unicode"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestRequestWorker(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO public\."soccer"`).WithArgs(
		driver.Value("SOCCER"),
		driver.Value("0.244"),
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	SingleRequestWorker(3, db)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expections: %s", err)
	}
}

func SingleRequestWorker(timeoutInterval int, db *sql.DB) {
	requestedString := `{"lines":{"SOCCER":"0.244"}}`

	sportRatio := strings.TrimFunc(string(requestedString), func(r rune) bool {
		return !unicode.IsDigit(r)
	})

	db_storage.PutSportLine(db, "SOCCER", sportRatio)	// Сохраняем в БД
}