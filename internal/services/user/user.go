package user

import (
	"fmt"
	"strings"
	"time"

	models "github.com/MomentsFromEarth/api/internal/models/user"
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

func formatNewUser(newUser *models.NewUser) *models.User {
	now := genTimestamp()
	userID, _ := genUserID()
	username := userID
	return &models.User{
		MFEKey:          getMfeKey(userID),
		UserID:          userID,
		UserName:        username,
		Email:           newUser.Email,
		Avatar:          newUser.Avatar,
		CognitoSub:      newUser.CognitoSub,
		JoinMailingList: newUser.JoinMailingList,
		NewUser:         true,
		Created:         now,
		Updated:         now,
		QueryKey01:      getEmailQueryKey(newUser.Email),
		QueryKey02:      getUserNameQueryKey(username),
	}
}

func formatEmailQueryInput(email string) *dynamodb.QueryInput {
	return &dynamodb.QueryInput{
		TableName:              aws.String(mfeTableName),
		IndexName:              aws.String(fmt.Sprintf("%s-index", mfeQuery01)),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :qk01", mfeQuery01)),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":qk01": {
				S: aws.String(getEmailQueryKey(email)),
			},
		},
	}
}

func formatUserNameQueryInput(username string) *dynamodb.QueryInput {
	return &dynamodb.QueryInput{
		TableName:              aws.String(mfeTableName),
		IndexName:              aws.String(fmt.Sprintf("%s-index", mfeQuery02)),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :qk02", mfeQuery02)),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":qk02": {
				S: aws.String(getUserNameQueryKey(username)),
			},
		},
	}
}

func getMfeKey(userID string) string {
	return fmt.Sprintf("usr:%s", userID)
}

func getUserNameQueryKey(username string) string {
	// todo: strip out any unwanted characters
	return fmt.Sprintf("usr:username:%s", strings.ToLower(username))
}

func getEmailQueryKey(email string) string {
	return fmt.Sprintf("usr:email:%s", strings.ToLower(email))
}

// Create is a function
func Create(newUser *models.NewUser) (*models.User, error) {

	existingUser, err := Read(newUser.Email)
	if existingUser != nil {
		return existingUser, nil
	}

	user := formatNewUser(newUser)

	userInput, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      userInput,
		TableName: aws.String(mfeTableName),
	}

	_, err = dynamoClient().PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return user, nil
}

// Read is a function
func Read(email string) (*models.User, error) {
	dydb := dynamoClient()
	result, err := dydb.Query(formatEmailQueryInput(email))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if len(result.Items) == 0 {
		msg := fmt.Errorf("User not found: %s", email)
		return nil, msg
	}

	res := result.Items[0]
	user := &models.User{}
	err = dynamodbattribute.UnmarshalMap(res, &user)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return user, nil
}

// Update is a function
func Update(user *models.User) (*models.User, error) {

	existingUser, err := Read(user.Email)
	if err != nil {
		return nil, err
	}

	user.MFEKey = existingUser.MFEKey
	user.Email = existingUser.Email
	user.QueryKey01 = existingUser.QueryKey01
	user.QueryKey02 = getUserNameQueryKey(user.UserName)
	user.CognitoSub = existingUser.CognitoSub
	user.UserID = existingUser.UserID
	user.Created = existingUser.Created
	user.Updated = genTimestamp()

	userInput, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      userInput,
		TableName: aws.String(mfeTableName),
	}
	_, err = dynamoClient().PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return user, nil
}

// Delete is a function
func Delete(email string) (*models.User, error) {
	existingUser, err := Read(email)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"mfe_key": {
				S: aws.String(existingUser.MFEKey),
			},
		},
		TableName: aws.String(mfeTableName),
	}

	_, err = dynamoClient().DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return existingUser, nil
}

// Profile is a function
func Profile(username string) (*models.ProfileUser, error) {
	dydb := dynamoClient()
	result, err := dydb.Query(formatUserNameQueryInput(username))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if len(result.Items) == 0 {
		msg := fmt.Errorf("User not found: %s", username)
		return nil, msg
	}

	res := result.Items[0]
	user := &models.ProfileUser{}
	err = dynamodbattribute.UnmarshalMap(res, &user)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return user, nil
}
