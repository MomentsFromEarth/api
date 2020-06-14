package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
)

type token struct {
	Email string `json:"email"`
}

// Init is the entrypoint of api package
func Init() *gin.Engine {
	r := gin.Default()
	r.Use(authenticate)
	r.GET("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "user",
		})
	})
	return r
}

func authenticate(c *gin.Context) {

	keySet, err := jwk.Fetch("https://cognito-idp.us-east-1.amazonaws.com/us-east-1_k0qjPXdf5/.well-known/jwks.json")
	if err != nil {
		log.Printf("failed to parse JWK: %s", err)
		return
	}

	authToken := c.GetHeader("Authorization")
	splitToken := strings.Split(authToken, "Bearer ")
	authToken = splitToken[1]
	authToken = "eyJraWQiOiI5ckF2Tnh3bG55N2hxcWRuN3dMVlhSTklqK25zMXd3Sk53M3hQcVZPWHJrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiIxMDJkYTEzNy03MzE1LTRlOTctYThiYy1iZGZkZWRiMDgxNDYiLCJhdWQiOiJlNXJnYTE5N2ZoaXRlbWoyZmZzNWdvMXUyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV2ZW50X2lkIjoiOWE2OWQ5OTUtOWUwMS00Y2YxLWE3ZDUtMjE1MTdhY2Y2NTZhIiwidG9rZW5fdXNlIjoiaWQiLCJhdXRoX3RpbWUiOjE1OTIwOTE0NDAsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC51cy1lYXN0LTEuYW1hem9uYXdzLmNvbVwvdXMtZWFzdC0xX2swcWpQWGRmNSIsImNvZ25pdG86dXNlcm5hbWUiOiJ3aWxsc3RlcHBAZ21haWwuY29tIiwiZXhwIjoxNTkyMDk1MDQwLCJpYXQiOjE1OTIwOTE0NDAsImVtYWlsIjoid2lsbHN0ZXBwQGdtYWlsLmNvbSJ9.PnpqXm6oXYjz46hiYXggL-ESAldc-xfMQHLfrKCX2GWeCCivr6nO8O1H2zQqZcRkYYbV-pPSGkRAfPqloNqVgLqqKBGBQg3F3vGz467wgkSSBX3OSVnZUSqfFLxhd6gA4LLn3Xa31pKkwfR21LL4JO4xpR6BNMAbhJRYOKT2e2LJs73HX7Bx-LVg_hBRSybgCqRWQAbLAU3yCAzPGgZ9o1vYNrbocUGXZQgbQGFCevqTTFdYNLEv_BAsmYE_hmx4Dq6j_8chvy40tLXR-rl3UBQA6q9HP1NdODMmLimgmwpWzct8i1fidgs0Jt9dTwxmV5cyAp6SFnhB2tc-nHYr1Q"

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
					err = errors.New("token expired")
					log.Println("token expired")
				}
			} else {
				log.Println("Invalid claims for id token")
				log.Println(err)
			}

			j, _ := json.Marshal(claims)
			t := token{}
			json.Unmarshal([]byte(j), &t)
			fmt.Println("This is my email")
			fmt.Println(t.Email)
		}
	} else {
		log.Println("Invalid")
		res2B, _ := json.Marshal(parsed)
		fmt.Println(string(res2B))
	}
	// https://cognito-idp.{region}.amazonaws.com/{userPoolId}/.well-known/jwks.json
	/*
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRS256); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
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
		if err == nil && token.Valid {
			fmt.Println("Your token is valid.  I like your style.")
		} else {
			fmt.Println("This token is terrible!  I cannot accept this.")
		}
	*/
	c.Next()
}

func init() {
	fmt.Println("api.init")
}
