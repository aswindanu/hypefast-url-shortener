package urlshortener

import (
	"encoding/json"
	"net/http"
	"strconv"

	"hypefast-url-shortener/app/handler"
	"hypefast-url-shortener/app/model"
	"hypefast-url-shortener/util"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func GetAllUrls(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	urls := []model.Url{}
	db.Find(&urls)
	handler.RespondJSON(w, http.StatusOK, urls)
}

func CreateUrl(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	url := model.Url{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&url); err != nil {
		handler.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	url.Uid = util.RandStringBytes(6)
	url.AlternateURL = `https://tinyurl.com/` + url.Uid
	url.IPOrigin = util.ReadUserIP(r)
	url.Activate()
	if err := db.Save(&url).Error; err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.RespondJSON(w, http.StatusCreated, url)
}

func GetUrl(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uid := vars["uid"]
	url := GetUrlOr404(db, uid, w, r)
	if url == nil {
		return
	}

	if err := db.Save(&url).Error; err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !url.IsActive {
		handler.RespondError(w, http.StatusBadRequest, "Invalid is_active status: "+strconv.FormatBool(url.IsActive))
		return
	}

	resp := make(map[string]interface{})
	resp["redirect_to"] = url.OriginalURL
	resp["redirect_count"] = url.RedirectCount
	resp["created_at"] = url.CreatedAt.String()
	handler.RespondJSON(w, http.StatusOK, resp)
}

func RedirectUrl(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	alternateUrl := r.URL.Query().Get("alternate_url")

	uid := alternateUrl[len(alternateUrl)-6:]
	url := GetUrlOr404(db, uid, w, r)
	if url == nil {
		return
	}

	if !url.IsActive {
		handler.RespondError(w, http.StatusBadRequest, "Invalid is_active status: "+strconv.FormatBool(url.IsActive))
		return
	}

	url.RedirectCount += 1
	if err := db.Save(&url).Error; err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
}

func ActivateUrl(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uid := vars["uid"]
	url := GetUrlOr404(db, uid, w, r)
	if url == nil {
		return
	}
	url.Activate()
	if err := db.Save(&url).Error; err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.RespondJSON(w, http.StatusOK, url)
}

func DeactivateUrl(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uid := vars["uid"]
	url := GetUrlOr404(db, uid, w, r)
	if url == nil {
		return
	}
	url.Deactivate()
	if err := db.Save(&url).Error; err != nil {
		handler.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.RespondJSON(w, http.StatusNoContent, nil)
}

// GetUrlOr404 gets a project instance if exists, or respond the 404 error otherwise
func GetUrlOr404(db *gorm.DB, uid string, w http.ResponseWriter, r *http.Request) *model.Url {
	url := model.Url{}
	if err := db.First(&url, model.Url{Uid: uid}).Error; err != nil {
		handler.RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &url
}
