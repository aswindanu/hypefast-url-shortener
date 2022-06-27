package urlshortener_test

import (
	// "net/http"
	// "io/ioutil"

	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	// app

	// handler
	"hypefast-url-shortener/app/handler/urlshortener"

	// config
	"hypefast-url-shortener/config"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func GetDB() *gorm.DB {
	var err error
	var db *gorm.DB
	var dbURI string

	_conf := config.GetConfig()

	// MYSQL
	if _conf.DB.Dialect == "mysql" {
		dbURI = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
			_conf.DB.Username,
			_conf.DB.Password,
			_conf.DB.Host,
			_conf.DB.Port,
			_conf.DB.Name,
			_conf.DB.Charset,
		)
	}
	// POSTGRES
	if _conf.DB.Dialect == "postgres" {
		dbURI = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			// _conf.DB.Host,
			"127.0.0.1",
			_conf.DB.Username,
			_conf.DB.Password,
			_conf.DB.Name,
			// _conf.DB.Port,
			5454,
		)
	}

	db, err = gorm.Open(_conf.DB.Dialect, dbURI)
	if err != nil {
		fmt.Println("dbURI: ", dbURI)
		log.Fatal("Could not connect database:" + err.Error())
	}
	return db
}

func TestGetUrlOr404(t *testing.T) {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	db := GetDB()
	res := urlshortener.GetUrlOr404(db, "none01", w, r)

	if res != nil {
		t.Errorf("expected test nil got %v", res.Uid)
	}

	res = urlshortener.GetUrlOr404(db, "fixed1", w, r)
	if res.Uid != "fixed1" {
		t.Errorf("expected test 'fixed1' got %v", res.Uid)
	}
}

func TestGetUrls(t *testing.T) {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)

	r, _ := http.NewRequest(http.MethodGet, "/api/v1/url-shortener", nil)
	w := httptest.NewRecorder()
	db := GetDB()

	urlshortener.GetAllUrls(db, w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUrl(t *testing.T) {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)

	// 200
	r, _ := http.NewRequest(http.MethodGet, "/api/v1/url-shortener/fixed1", nil)
	w := httptest.NewRecorder()
	db := GetDB()
	vars := map[string]string{
		"uid": "fixed1",
	}
	r = mux.SetURLVars(r, vars)
	urlshortener.GetUrl(db, w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	// 404
	r, _ = http.NewRequest(http.MethodGet, "/api/v1/url-shortener/none01", nil)
	w = httptest.NewRecorder()
	db = GetDB()
	vars = map[string]string{
		"uid": "none01",
	}
	r = mux.SetURLVars(r, vars)
	urlshortener.GetUrl(db, w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// 400
	r, _ = http.NewRequest(http.MethodGet, "/api/v1/url-shortener/deact1", nil)
	w = httptest.NewRecorder()
	db = GetDB()
	vars = map[string]string{
		"uid": "deact1",
	}
	r = mux.SetURLVars(r, vars)
	urlshortener.GetUrl(db, w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUrl(t *testing.T) {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)

	// 200
	jsonStr := []byte(`{"original_url":"https://www.google.com"}`)
	r, _ := http.NewRequest(http.MethodPost, "/api/v1/url-shortener", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	db := GetDB()
	urlshortener.CreateUrl(db, w, r)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 500
	jsonStr = []byte(`{"wrong_key":"https://www.google.com"}`)
	r, _ = http.NewRequest(http.MethodPost, "/api/v1/url-shortener", bytes.NewBuffer(jsonStr))
	w = httptest.NewRecorder()
	db = GetDB()
	urlshortener.CreateUrl(db, w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 500
	jsonStr = []byte(`{"original_url":123456}`)
	r, _ = http.NewRequest(http.MethodPost, "/api/v1/url-shortener", bytes.NewBuffer(jsonStr))
	w = httptest.NewRecorder()
	db = GetDB()
	urlshortener.CreateUrl(db, w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestARedirectUrl(t *testing.T) {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)

	// 301
	r, _ := http.NewRequest(http.MethodGet, "/api/v1/url-shortener/redirect?alternate_url=https://tinyurl.com/fixed1", nil)
	w := httptest.NewRecorder()
	db := GetDB()
	urlshortener.RedirectUrl(db, w, r)
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	// 400
	r, _ = http.NewRequest(http.MethodGet, "/api/v1/url-shortener/redirect?alternate_url=https://tinyurl.com/deact1", nil)
	w = httptest.NewRecorder()
	db = GetDB()
	urlshortener.RedirectUrl(db, w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestActivateUrl(t *testing.T) {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)

	r, _ := http.NewRequest(http.MethodPost, "/api/v1/url-shortener/fixed1/activate", nil)
	w := httptest.NewRecorder()
	db := GetDB()
	urlshortener.ActivateUrl(db, w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeactivateUrl(t *testing.T) {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)

	r, _ := http.NewRequest(http.MethodPost, "/api/v1/url-shortener/fixed1/deactivate", nil)
	w := httptest.NewRecorder()
	db := GetDB()
	urlshortener.DeactivateUrl(db, w, r)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
