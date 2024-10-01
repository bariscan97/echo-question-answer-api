package utils

import (
	"fmt"
	"os"
	"articles-api/models/token"
	"time"
	"articles-api/models/user"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJwtToken(user *user.FetchUserModel) (string, error){
	
	expirationTime := time.Now().Add(72 * time.Hour)
	
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))
	
	claims := &token.Claim{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return tokenString, nil
}