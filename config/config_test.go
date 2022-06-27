package config_test

import (
	// "net/http"
	// "io/ioutil"

	"fmt"
	"log"
	"os"
	"testing"

	"hypefast-url-shortener/config"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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

func TestConfig(t *testing.T) {
	GetDB()
	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	// _dbconf := config.DBConfig{}
	// _conf := config.Config{}
	// // db := GetDB()

	// if _conf.DB == nil {
	// 	t.Errorf("expected test got %v", nil)
	// }
	// if _dbconf.Name != "root" {
	// 	t.Errorf("expected test 'root' got %v", nil)
	// }
	// if db.Error != nil {
	// 	t.Errorf("expected test 'db' got %v", db.Error)
	// }
}
