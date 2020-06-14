package dynamodbmfe

import (
	awsmfe "github.com/MomentsFromEarth/api/internal/services/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Client returns AWS DynamoDB client
func Client() *dynamodb.DynamoDB {
	sess, _ := awsmfe.Session()
	return dynamodb.New(sess)
}
