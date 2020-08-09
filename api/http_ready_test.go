package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckConnection(t *testing.T) {
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