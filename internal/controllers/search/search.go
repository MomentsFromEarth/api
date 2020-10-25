package search

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "strings"
  "time"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/signer/v4"
  "github.com/gin-gonic/gin"
)


// Query is an endpoint
func Query(c *gin.Context) {
	  // Basic information for the Amazon Elasticsearch Service domain
	  domain := "https://search-mfe-dzlkaaoblbbhhxcqd6shu4oo3u.us-east-1.es.amazonaws.com" // e.g. https://my-domain.region.es.amazonaws.com
	  // index := "my-index"
	  // id := "1"
	  endpoint := domain + "/_stats"
	  region := "us-east-1" // e.g. us-east-1
	  service := "es"
	
	  // Sample JSON document to be included as the request body
	  json := `{ "title": "Thor: Ragnarok", "director": "Taika Waititi", "year": "2017" }`
	  body := strings.NewReader(json)
	
	  // Get credentials from environment variables and create the AWS Signature Version 4 signer
	  credentials := credentials.NewEnvCredentials()
	  fmt.Print("CREDS")
	  fmt.Print(credentials)
	  signer := v4.NewSigner(credentials)
	
	  // An HTTP client for sending the request
	  client := &http.Client{}
	
	  // Form the HTTP request
	  req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	  if err != nil {
		fmt.Print(err)
	  }
	
	  // You can probably infer Content-Type programmatically, but here, we just say that it's JSON
	  req.Header.Add("Content-Type", "application/json")
	
	  // Sign the request, send it, and print the response
	  signer.Sign(req, body, service, region, time.Now())
	  resp, err := client.Do(req)
	  if err != nil {
		fmt.Print(err)
	  }
	  fmt.Print(resp.Status + "\n")
	  respbody, err := ioutil.ReadAll(resp.Body)
	  resp.Body.Close()
	  fmt.Print(string(respbody))

	  c.JSON(http.StatusOK, gin.H{"message": "ok"})
	  c.Next()
}
