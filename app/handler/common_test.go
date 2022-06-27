package handler_test

import (
	// "net/http"
	// "io/ioutil"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"hypefast-url-shortener/app/handler"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGetRespondJSON(t *testing.T) {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	type Payload struct {
		Payload string
	}
	payload := Payload{}

	w := httptest.NewRecorder()
	handler.RespondJSON(w, 200, map[string]string{"payload": "test"})

	res := w.Result()
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&payload); err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	defer res.Body.Close()

	if payload.Payload != "test" {
		t.Errorf("expected test got %v", payload.Payload)
	}
}

func TestGetRespondError(t *testing.T) {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	type Error struct {
		Error string
	}
	errResponse := Error{}

	w := httptest.NewRecorder()
	handler.RespondError(w, 500, "test")

	res := w.Result()
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&errResponse); err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	defer res.Body.Close()

	if errResponse.Error != "test" {
		t.Errorf("expected test got %v", errResponse.Error)
	}
}
