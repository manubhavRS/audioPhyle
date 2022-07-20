package utilities

import (
	"audioPhile/models"
	"cloud.google.com/go/firestore"
	cloud "cloud.google.com/go/storage"
	"database/sql"
	"encoding/json"
	firebase "firebase.google.com/go"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const ContextUserKey string = "user"
const TokenString = "Token"
const ShippingCharges = 50

func FetchExpectDateOfDelivery() time.Time {
	return time.Now().Add(time.Hour * 120)
}

func FetchExpireTime() time.Time {
	return time.Now().Add(time.Hour * 30)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func JsonData(x interface{}) ([]byte, error) {
	jsonResponse, err := json.Marshal(x)
	return jsonResponse, err
}

func CreateImageUrl(imagePath string) models.ImageStructure {
	bucket := os.Getenv("bucketName")
	imageStructure := models.ImageStructure{
		ImageName: imagePath,
		URL:       "https://storage.cloud.google.com/" + bucket + "/image/" + imagePath,
	}
	return imageStructure
}

type App struct {
	Ctx     context.Context
	Client  *firestore.Client
	Storage *cloud.Client
}

func FireBaseUpload(file multipart.File, imagePath string) error {
	route := App{}
	route.Ctx = context.Background()
	serviceKey := os.Getenv("serviceKey")
	sa := option.WithCredentialsJSON([]byte(serviceKey))
	app, err := firebase.NewApp(route.Ctx, nil, sa)
	if err != nil {
		return err
	}
	route.Client, err = app.Firestore(route.Ctx)
	if err != nil {
		return err
	}
	route.Storage, err = cloud.NewClient(route.Ctx, sa)
	if err != nil {
		return err
	}
	bucket := os.Getenv("bucketName")
	/*	_, _, err = route.Client.Collection(bucket).Add(route.Ctx, "image")
		if err != nil {
			return err
		}
	*/
	wc := route.Storage.Bucket(bucket).Object("image/" + imagePath).NewWriter(route.Ctx)
	_, err = io.Copy(wc, file)
	if err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}
func TokenResponse(tokenString string) ([]byte, error) {
	resp := make(map[string]string)
	resp["Name"] = "token"
	resp["Value"] = tokenString
	resp["Expires"] = FetchExpireTime().String()
	jsonResponse, err := JsonData(resp)
	return jsonResponse, err
}
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func ArgumentsMapping(r *http.Request) (models.ProductListSearchModel, error, int) {
	var productSearch models.ProductListSearchModel
	pageNo := r.URL.Query().Get("pageNo")
	if r.URL.Query().Get("category") != "" {
		productSearch.Category = strings.Split(r.URL.Query().Get("category"), ",")
	}
	productSearch.Search = r.URL.Query().Get("search")
	limit := r.URL.Query().Get("limit")
	var err error
	if len(r.URL.Query().Get("latest")) > 0 {
		productSearch.Latest, err = strconv.ParseBool(r.URL.Query().Get("latest"))
		if err != nil {
			return productSearch, err, http.StatusBadRequest
		}
	}
	if len(r.URL.Query().Get("youMayLike")) > 0 {
		productSearch.YML, err = strconv.ParseBool(r.URL.Query().Get("youMayLike"))
		if err != nil {
			return productSearch, err, http.StatusBadRequest
		}
	}
	if len(productSearch.Category) == 0 {
		productSearch.CheckCategory = true
	} else {
		productSearch.CheckCategory = false
	}
	if productSearch.Latest == true && productSearch.YML == true {
		return productSearch, err, http.StatusBadRequest
	}
	if productSearch.Latest == true {
		productSearch.Limit = 1
		productSearch.OrderBy = "created_by DESC"
		productSearch.PageNo = 1
	} else if productSearch.YML == true {
		productSearch.Limit = 5
		productSearch.OrderBy = "RANDOM ()"
		productSearch.PageNo = 1
	} else {
		productSearch.PageNo, _ = strconv.Atoi(pageNo)
		productSearch.Limit, _ = strconv.Atoi(limit)
	}
	return productSearch, nil, http.StatusOK
}
