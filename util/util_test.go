package util_test

import (
	// "net/http"
	// "io/ioutil"

	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"hypefast-url-shortener/config"
	"hypefast-url-shortener/util"

	"github.com/jinzhu/gorm"
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

func TestGenerateUID(t *testing.T) {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	util.GenerateID()
	util.RandStringBytes(6)
}

func TestGetIP(t *testing.T) {
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	util.ReadUserIP(r)
}
