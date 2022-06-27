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
	db.DropTableIfExists(&Url{})
	db.AutoMigrate(&Url{})

	// Add seed
	db.Exec(`
	INSERT INTO "public"."urls" ("id", "created_at", "updated_at", "deleted_at", "uid", "original_url", "alternate_url", "redirect_count", "ip_origin", "is_active") VALUES ('1', '2022-06-27 15:50:20.256583+07', '2022-06-27 16:14:48.234717+07', NULL, 'deact1', 'https://google.com', 'https://tinyurl.com/123456', '5', '172.20.0.1:60344', 'f');
	INSERT INTO "public"."urls" ("id", "created_at", "updated_at", "deleted_at", "uid", "original_url", "alternate_url", "redirect_count", "ip_origin", "is_active") VALUES ('2', '2022-06-27 15:50:20.256583+07', '2022-06-27 16:14:48.234717+07', NULL, 'fixed1', 'https://google.com', 'https://tinyurl.com/fixed1', '5', '172.20.0.1:60344', 't');
	SELECT setval('urls_id_seq', 3, true);`)
	return db
}
