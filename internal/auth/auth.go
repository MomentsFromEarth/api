package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
)

var keySet *jwk.Set

type token struct {
	Email string `json:"email"`
}

// Init is the entrypoint of auth package
func Init() {
	var err error
	keySet, err = jwk.Fetch("https://cognito-idp.us-east-1.amazonaws.com/us-east-1_iJhXcJLvi/.well-known/jwks.json")
	if err != nil {
		log.Printf("Failed to parse JWK: %s", err)
	}
}

// Valid is a function
func Valid(key string) bool {
	apiKey := os.Getenv("API_KEY")
	return apiKey == key
}

// Run is the entrypoint of auth package
func Run() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			unauthorized(c, "Authorization header missing")
		}

		splitToken := strings.Split(authToken, "Bearer ")
		authToken = splitToken[1]

		parsed, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			kid, ok := token.Header["kid"].(string)
			if !ok {
				unauthorized(c, "Token KID header not found")
			}
			keys := keySet.LookupKeyID(kid)
			if len(keys) == 0 {
				unauthorized(c, "Token Key not found")
			}

			var raw interface{}
			return raw, keys[0].Raw(&raw)
		})

		if err == nil && parsed.Valid {
			if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
				err = claims.Valid()
				if err != nil {
					unauthorized(c, "Token contains invalid data")
				}
				j, _ := json.Marshal(claims)
				t := token{}
				json.Unmarshal([]byte(j), &t)
				email := t.Email
				if empty(email) {
					unauthorized(c, "Email not found, cannot continue")
				}
				fmt.Printf("[AuthSuccess] %v", t.Email)
				c.Set("email", t.Email)
			}
		} else {
			unauthorized(c, "Token invalid")
		}
		c.Next()
	}
}

func empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func unauthorized(c *gin.Context, message string) {
	fmt.Printf("[AuthError] %v", message)
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": message})
}
