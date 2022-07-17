package utilities

import (
	"audioPhile/models"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

const ContextUserKey string = "user"
const TokenString = "Token"
const ShippingCharges = 50

func FetchExpectDateOfDelivery() time.Time {
	return time.Now().Add(time.Hour * 120)
}

func FetchExpireTime() time.Time {
	return time.Now().Add(time.Minute * 30)
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
