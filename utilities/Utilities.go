package utilities

import (
	"golang.org/x/crypto/bcrypt"
)

const Secretkey string = "SuperSecretKey"
const ContextUserKey string = "user"
const TokenString = "Token"
const ContextRefreshToken string = "refreshToken"
const AdminRole = "admin"
const SubAdminRole = "sub-admin"
const UserRole = "user"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
