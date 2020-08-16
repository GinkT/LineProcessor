package api

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckConnection(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	StatusCheckInit(db, "")
	if err = db.Ping(); err != nil {
		t.Fatal("ggwp db ping not gone")
	}

	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		logrus.Fatalln(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CheckConnection)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `DB CONNECTED`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}