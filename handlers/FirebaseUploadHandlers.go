package handlers

import (
	"audioPhile/database"
	"audioPhile/database/helper"
	"audioPhile/models"
	_ "audioPhile/models"
	"cloud.google.com/go/firestore"
	cloud "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	_ "firebase.google.com/go/v4/auth"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
	ctx     context.Context
	client  *firestore.Client
	storage *cloud.Client
}

func FireBaseUploadHandler(w http.ResponseWriter, r *http.Request) {
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		var err error
		imagePath := "image_" + time.Now().String() + ".jpg"
		ps := models.ProductAssetModel{
			ID:        r.URL.Query().Get("productID"),
			AssetName: imagePath,
		}
		_, err = helper.AddProductAssetHelper(ps, tx)
		if err != nil {
			return err
		}
		route := App{}
		route.ctx = context.Background()
		serviceKey := os.Getenv("serviceKey")
		sa := option.WithCredentialsJSON([]byte(serviceKey))
		app, err := firebase.NewApp(route.ctx, nil, sa)
		if err != nil {
			return err
		}
		route.client, err = app.Firestore(route.ctx)
		if err != nil {
			return err
		}
		route.storage, err = cloud.NewClient(route.ctx, sa)
		if err != nil {
			return err
		}
		bucket := os.Getenv("bucketName")
		route.client.Collection(bucket).Add(route.ctx, "image")
		file, _, err := r.FormFile("image")
		r.ParseMultipartForm(10 << 20)
		if err != nil {
			return err
		}
		defer file.Close()
		log.Printf(imagePath)
		wc := route.storage.Bucket(bucket).Object("image/" + imagePath).NewWriter(route.ctx)
		_, err = io.Copy(wc, file)
		if err != nil {
			return err
		}
		if err := wc.Close(); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		log.Printf("FireBaseUploadHandler: %v", txErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
