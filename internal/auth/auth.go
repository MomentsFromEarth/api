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
