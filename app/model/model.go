package model

import (
	// "time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Url struct {
	gorm.Model
	Uid           string `json:"uid" gorm:"unique; type:varchar(25); default:null"`
	OriginalURL   string `json:"original_url" gorm:"not null; type:varchar(500); default:null"`
	AlternateURL  string `json:"alternate_url" gorm:"type:varchar(100); default:null"`
	RedirectCount int    `json:"redirect_count" gorm:"type:int; default:0"`
	IPOrigin      string `json:"ip_origin" gorm:"type:varchar(20); default:null"`
	IsActive      bool   `json:"is_active" gorm:"type:bool; default:true"`
}

func (u *Url) Activate() {
	u.IsActive = true
}

func (u *Url) Deactivate() {
	u.IsActive = false
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Url{})
	return db
}
