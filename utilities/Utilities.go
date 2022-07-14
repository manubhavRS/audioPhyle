package utilities

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

const Secretkey string = "SuperSecretKey"
const ContextUserKey string = "user"
const TokenString = "Token"
const ContextRefreshToken string = "refreshToken"
const AdminRole = "admin"
const SubAdminRole = "sub-admin"
const UserRole = "user"
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
