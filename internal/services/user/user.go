package user

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MomentsFromEarth/api/internal/models"
	dynamodbmfe "github.com/MomentsFromEarth/api/internal/services/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/teris-io/shortid"
)

var mfeTableName = "MFE"
var mfeQuery01 = "query_key_01"
var mfeQuery02 = "query_key_02"

func dynamoClient() *dynamodb.DynamoDB {
	return dynamodbmfe.Client()
}

func genTimestamp() int64 {
	return time.Now().Unix()
}

func genUserID() (string, error) {
	return shortid.Generate()
}

func formatUser(newUser *models.NewUser) *models.User {
	now := genTimestamp()
	userID, _ := genUserID()
	username := userID
	return &models.User{
		WSKey:           fmt.Sprintf("usr:%s", userID),
		UserID:          userID,
		UserName:        username,
		Email:           newUser.Email,
		Avatar:          newUser.Avatar,
		CognitoSub:      newUser.CognitoSub,
		JoinMailingList: newUser.JoinMailingList,
		NewUser:         true,
		Created:         now,
		Updated:         now,
		QueryKey01:      fmt.Sprintf("usr:email:%s", strings.ToLower(newUser.Email)),
		QueryKey02:      fmt.Sprintf("usr:username:%s", strings.ToLower(newUser.Email)),
	}
}

func fromEmailParams(email string) *dynamodb.QueryInput {
	return &dynamodb.QueryInput{
		TableName:              aws.String(mfeTableName),
		IndexName:              aws.String(fmt.Sprintf("%s-index", mfeQuery01)),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :qk01", mfeQuery01)),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":qk01": {
				S: aws.String(fmt.Sprintf("usr:email:%s", strings.ToLower(email))),
			},
		},
	}
}

// Get is a function
func Get(email string) *models.User {
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
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return user
}

// Create is a function
func Create(newUser *models.NewUser) *models.User {
	user := formatUser(newUser)
	return user
}
