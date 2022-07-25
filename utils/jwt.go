package utils

import (
   "fmt"
   "github.com/dgrijalva/jwt-go"
   "github.com/gin-gonic/gin"
)

var SECRET_KEY = []byte("gosecretkey")

func CreateJwtToken(c *gin.Context, userId string) (string,error) {
    type MyCustomClaims struct {
        ID string `json:"userId"`
        jwt.StandardClaims
    }

    // Create the Claims
    claims := MyCustomClaims{
        userId,
        jwt.StandardClaims{
            ExpiresAt: 15000,
            Issuer:    "test",
        },
    }

    token:= jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(SECRET_KEY)

    fmt.Println(tokenString)

    return tokenString, err
}