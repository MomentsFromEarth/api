package user

import (
	"fmt"
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

func formatNewUser(newUser *models.NewUser) *models.User {
	now := genTimestamp()
	userID, _ := genUserID()
	username := userID
	return &models.User{
		MFEKey:          fmt.Sprintf("usr:%s", userID),
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

func formatEmailQueryInput(email string) *dynamodb.QueryInput {
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

func mapProfileToUser(profileUser *models.ProfileUser) *models.User {
	return &models.User{
		MFEKey:          profileUser.MFEKey,
		Email:           profileUser.Email,
		UserID:          profileUser.UserID,
		UserName:        profileUser.UserName,
		Avatar:          profileUser.Avatar,
		CognitoSub:      profileUser.CognitoSub,
		Created:         profileUser.Created,
		Updated:         profileUser.Updated,
		JoinMailingList: profileUser.JoinMailingList,
		NewUser:         profileUser.NewUser,
		QueryKey01:      profileUser.QueryKey01,
		QueryKey02:      profileUser.QueryKey02,
	}
}

func mapUserToProfile(user *models.User) *models.ProfileUser {
	return &models.ProfileUser{
		MFEKey:          user.MFEKey,
		Email:           user.Email,
		UserID:          user.UserID,
		UserName:        user.UserName,
		Avatar:          user.Avatar,
		CognitoSub:      user.CognitoSub,
		Created:         user.Created,
		Updated:         user.Updated,
		JoinMailingList: user.JoinMailingList,
		NewUser:         user.NewUser,
		QueryKey01:      user.QueryKey01,
		QueryKey02:      user.QueryKey02,
	}
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
