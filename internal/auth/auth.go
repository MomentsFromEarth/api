package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
)

type token struct {
	Email string `json:"email"`
}

// Run is the entrypoint of auth package
func Run(c *gin.Context) {
	keySet, err := jwk.Fetch("https://cognito-idp.us-east-1.amazonaws.com/us-east-1_k0qjPXdf5/.well-known/jwks.json")
	if err != nil {
		log.Printf("failed to parse JWK: %s", err)
		return
	}

	authToken := c.GetHeader("Authorization")
	splitToken := strings.Split(authToken, "Bearer ")
	authToken = splitToken[1]
	authToken = "eyJraWQiOiI5ckF2Tnh3bG55N2hxcWRuN3dMVlhSTklqK25zMXd3Sk53M3hQcVZPWHJrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiIxMDJkYTEzNy03MzE1LTRlOTctYThiYy1iZGZkZWRiMDgxNDYiLCJhdWQiOiJlNXJnYTE5N2ZoaXRlbWoyZmZzNWdvMXUyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV2ZW50X2lkIjoiMDI3NjNmYmYtMWYwYS00ZDFhLWEyMGUtNmNjOTllMDU3MzA5IiwidG9rZW5fdXNlIjoiaWQiLCJhdXRoX3RpbWUiOjE1OTIxMDYzNjgsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC51cy1lYXN0LTEuYW1hem9uYXdzLmNvbVwvdXMtZWFzdC0xX2swcWpQWGRmNSIsImNvZ25pdG86dXNlcm5hbWUiOiJ3aWxsc3RlcHBAZ21haWwuY29tIiwiZXhwIjoxNTkyMTA5OTY4LCJpYXQiOjE1OTIxMDYzNjgsImVtYWlsIjoid2lsbHN0ZXBwQGdtYWlsLmNvbSJ9.HedlaxJ2qPSYWJZj07XIV4esCmLBKSi91ULgjvzjG9mKYIjqF1fQtVxpYPZ5ijADALD7QSXkn3UkZnCF1jxM5JEKGjEE_2NuT5V8wcaWSY27LAiVGUAE2qSto87Trle2JfKaUDqkGTSPyFtU9s52zCylX7Gah7R-_5QS6YKcYRIwsUr2krP-C4liPUG_DIzmK1JxwZET8slCrayoVwiPsdRfqUXZYzcki1YvdbgGgvVKm39KgARiXZgc3SBNfSgJgvQ_-G5cjtYd0qxlwAi86oL59-DTSs2Pf7T7ajDseu_dyvcNof8N6pT8D_hTp4jRIbbj_nhai0GY1Qpb51mpLg"

	parsed, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}
		keys := keySet.LookupKeyID(kid)
		if len(keys) == 0 {
			return nil, fmt.Errorf("key %v not found", kid)
		}

		var raw interface{}
		return raw, keys[0].Raw(&raw)
	})

	if err == nil && parsed.Valid {
		if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
			err = claims.Valid()
			if err == nil {
				now := time.Now().Unix()
				if claims.VerifyExpiresAt(now, true) == false {
					abort(c, "Token expired")
				}
			} else {
				abort(c, "Token contains invalid data")
			}
			j, _ := json.Marshal(claims)
			t := token{}
			json.Unmarshal([]byte(j), &t)
			fmt.Println(t.Email)
			c.Set("email", t.Email)
		}
	} else {
		abort(c, "Invalid Token")
	}
	c.Next()
}

func abort(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": message})
}
