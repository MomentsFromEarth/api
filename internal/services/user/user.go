package user

import (
	"fmt"
	"os"

	"github.com/MomentsFromEarth/api/internal/models"
	dynamodbmfe "github.com/MomentsFromEarth/api/internal/services/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var mfeTableName = "MFE"
var mfeQuery01 = "query_key_01"

func dynamoClient() *dynamodb.DynamoDB {
	return dynamodbmfe.Client()
}

func fromEmailParams(email string) *dynamodb.QueryInput {
	return &dynamodb.QueryInput{
		TableName:              aws.String(mfeTableName),
		IndexName:              aws.String(fmt.Sprintf("%s-index", mfeQuery01)),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :qk01", mfeQuery01)),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":qk01": {
				S: aws.String(fmt.Sprintf("usr:email:%s", email)),
			},
		},
	}
}

// FromEmail is a function
func FromEmail(email string) *models.User {
	dydb := dynamoClient()
	result, err := dydb.Query(fromEmailParams(email))
	if err != nil {
		fmt.Println(err.Error())
		return &models.User{}
	}
	res := result.Items[0]
	user := &models.User{}
	err = dynamodbattribute.UnmarshalMap(res, &user)
	if err != nil {
		fmt.Println("Got error unmarshalling:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return user
}
