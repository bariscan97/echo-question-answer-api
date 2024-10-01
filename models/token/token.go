package token

import (
	"github.com/golang-jwt/jwt/v4"
	"articles-api/models/user"
)

type Claim struct {	
	User  *user.FetchUserModel
	jwt.RegisteredClaims
}