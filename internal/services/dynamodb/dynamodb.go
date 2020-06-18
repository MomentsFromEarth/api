package dynamodbmfe

import (
	awsmfe "github.com/MomentsFromEarth/api/internal/services/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var client *dynamodb.DynamoDB

// Client returns AWS DynamoDB client
func Client() *dynamodb.DynamoDB {
	if client == nil {
		session := awsmfe.Session()
		client = dynamodb.New(session)
	}
	return client
}
