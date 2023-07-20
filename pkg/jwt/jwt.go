package jwt

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySecretString = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(id, usertype string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24 * 14).Unix()
	claims["id"] = id
	claims["usertype"] = usertype

	tokenString, err := token.SignedString(mySecretString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateTokenForPasswordReset(id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	claims["id"] = id

	tokenString, err := token.SignedString(mySecretString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateTokenForAccountVerification(id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	claims["id"] = id

	tokenString, err := token.SignedString(mySecretString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func IsTokenExpired(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySecretString, nil
	})
	if err != nil {
		return false
	}

	claims := token.Claims.(jwt.MapClaims)

	exp := claims["exp"].(float64)

	if time.Now().Unix() > int64(exp) {
		return true
	}

	return false
}

func IsValidToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySecretString, nil
	})
	if err != nil {
		return false
	}

	claims := token.Claims.(jwt.MapClaims)

	authorized := claims["authorized"].(bool)

	return authorized
}

func IsAdmin(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySecretString, nil
	})
	if err != nil {
		return false
	}

	claims := token.Claims.(jwt.MapClaims)

	userType := claims["usertype"].(string)

	if userType == "admin" {
		return true
	}

	return false
}

func IsHotelAdmin(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySecretString, nil
	})
	if err != nil {
		return false
	}

	claims := token.Claims.(jwt.MapClaims)

	userType := claims["usertype"].(string)

	if userType == "restaurant" {
		return true
	}

	return false
}

func GetUserID(tokenString string) string {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return mySecretString, nil
	})
	if err != nil {
		return ""
	}

	claims := token.Claims.(jwt.MapClaims)

	userID := claims["id"].(string)

	return userID
}
