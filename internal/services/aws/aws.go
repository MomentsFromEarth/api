package awsmfe

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var region = "us-east-1"

// Session returns new AWS Session
func Session() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
}
