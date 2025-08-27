package util

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pratyush934/tradealpha/server/models"
	"github.com/pratyush934/tradealpha/server/types"
)

var privateKey = []byte("iampratyushiamprayushiampratyushiampratyush")

/*
	1. CreateToken
	2. GetTokenFromHeader
	3. GetToken
*/

func CreateToken(u *models.User) (string, error) {

	ttl := 1800

	claims := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":    u.Id,
		"name":  u.Name,
		"email": u.Email,
		"role":  u.Role,
		"iat":   time.Now(),
		"eat":   time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	})

	return claims.SignedString(privateKey)
}

func GetToken(c echo.Context) (*jwt.Token, error) {
	header, err := GetTokenFromHeader(c)
	if err != nil {
		return nil, NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Not able to get TokenFromHeader jwt.go/GetToken", err)
	}

	parse, err := jwt.Parse(header, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return privateKey, nil
	})

	if err != nil {
		return nil, NewAppError(http.StatusInternalServerError, types.StatusInternalServerError, "Not able to parse the jwt", err)
	}
	return parse, nil

}

func GetTokenFromHeader(c echo.Context) (string, error) {
	str := c.Request().Header.Get("Authorization")
	newStr := strings.Split(str, " ")

	if len(newStr) != 2 || newStr[0] != "Bearer" {
		err := NewAppError(http.StatusBadGateway, types.StatusBadRequest, "Not able to get the newStr in jwt/GetTokenFromHeader", nil)
		return "", err
	}
	return newStr[1], nil

}
