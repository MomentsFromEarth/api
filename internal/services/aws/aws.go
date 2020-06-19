package awsmfe

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var region = "us-east-1"
var sess *session.Session = nil

// Session returns singleton AWS Session
func Session() *session.Session {
	if sess == nil {
		sess, _ = session.NewSession(&aws.Config{Region: aws.String(region)})
	}
	return sess
}

// CredSession returns new AWS Session
func CredSession(key string, secret string) *session.Session {
	return session.NewSession(&aws.Config{Region: aws.String(region)})
}
