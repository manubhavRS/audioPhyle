package handlers

import (
	"audioPhile/database"
	"audioPhile/database/helper"
	"audioPhile/models"
	_ "audioPhile/models"
	"audioPhile/utilities"
	_ "firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"time"
)

func FireBaseUploadHandler(w http.ResponseWriter, r *http.Request) {
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		var err error
		file, header, err := r.FormFile("image")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		ext, err := mime.ExtensionsByType(header.Header.Get("Content-Type"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}(file)
		imagePath := fmt.Sprintf("image-%v%v", time.Now().Unix(), ext[len(ext)-1])
		log.Printf(imagePath)

		ps := models.ProductAssetModel{
			ID:        r.URL.Query().Get("productID"),
			AssetName: imagePath,
		}
		_, err = helper.AddProductAssetHelper(ps, tx)
		if err != nil {
			return err
		}
		err = utilities.FireBaseUpload(file, imagePath)
		return err
	})
	if txErr != nil {
		log.Printf("FireBaseUploadHandler: %v", txErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
