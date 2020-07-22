package endpoints

import (
	jwt "github.com/dgrijalva/jwt-go"
	"strings"
)

type UserClaim struct {
	jwt.StandardClaims
	ID    string `json:"id"`
	Email string `json:"email"`
}

func GetTokenFromAuthorization(token string) string {
	return strings.Replace(token, "Bearer ", "", 1)
}

func DecodeUserFromToken(token string) (*UserClaim, error) {
	claim, err := jwt.ParseWithClaims(token, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("xxx"), nil
	})
	if err != nil {
		return nil, err
	}
	return claim.Claims.(*UserClaim), nil
}
