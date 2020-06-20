package moment

import (
	"fmt"
	"strings"
	"time"

	models "github.com/MomentsFromEarth/api/internal/models/moment"
	usermodels "github.com/MomentsFromEarth/api/internal/models/user"
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

func genMomentID() (string, error) {
	return shortid.Generate()
}

func genUserID() (string, error) {
	return shortid.Generate()
}

func formatNewMoment(newMoment *models.NewMoment) *models.Moment {
	now := genTimestamp()
	momentID, _ := genMomentID()
	status := "queued"
	return &models.Moment{
		MFEKey:      getMfeKey(momentID),
		MomentID:    momentID,
		Creator:     newMoment.Creator,
		Title:       newMoment.Title,
		Description: newMoment.Description,
		Filename:    newMoment.Filename,
		Type:        newMoment.Type,
		Size:        newMoment.Size,
		Status:      status,
		QueueID:     newMoment.QueueID,
		HostID:      "",
		Created:     now,
		Updated:     now,
		QueryKey01:  getStatusQueryKey(status),
		QueryKey02:  "mom",
	}
}

func formatNewUser(newUser *usermodels.NewUser) *usermodels.User {
	now := genTimestamp()
	userID, _ := genUserID()
	username := userID
	return &usermodels.User{
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

func mapProfileToUser(profileUser *usermodels.ProfileUser) *usermodels.User {
	return &usermodels.User{
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

func mapUserToProfile(user *usermodels.User) *usermodels.ProfileUser {
	return &usermodels.ProfileUser{
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

func getMfeKey(momentID string) string {
	return fmt.Sprintf("mom:%s", momentID)
}

func getStatusQueryKey(status string) string {
	return fmt.Sprintf("mom:status:%s", strings.ToLower(status))
}

func getUserNameQueryKey(username string) string {
	// todo: strip out any unwanted characters
	return fmt.Sprintf("usr:username:%s", strings.ToLower(username))
}

func getEmailQueryKey(email string) string {
	return fmt.Sprintf("usr:email:%s", strings.ToLower(email))
}

// Create is a function
func Create(newMoment *models.NewMoment) (*models.Moment, error) {
	moment := formatNewMoment(newMoment)
	momentInput, err := dynamodbattribute.MarshalMap(moment)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      momentInput,
		TableName: aws.String(mfeTableName),
	}
	_, err = dynamoClient().PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return moment, nil
}

// Read is a function
func Read(momentID string) (*models.Moment, error) {
	dydb := dynamoClient()
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"mfe_key": {
				S: aws.String(getMfeKey(momentID)),
			},
		},
		TableName: aws.String(mfeTableName),
	}
	result, err := dydb.GetItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	res := result.Item
	mom := &models.Moment{}
	err = dynamodbattribute.UnmarshalMap(res, &mom)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return mom, nil
}

// Update is a function
func Update(moment *models.Moment) (*models.Moment, error) {
	existing, err := Read(moment.MomentID)
	if err != nil {
		return nil, err
	}
	existing.Title = moment.Title
	existing.Description = moment.Description
	existing.Updated = genTimestamp()
	existing.HostID = moment.HostID
	existing.Status = strings.ToLower(moment.Status)
	existing.QueryKey01 = getStatusQueryKey(moment.Status)
	momentInput, err := dynamodbattribute.MarshalMap(existing)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      momentInput,
		TableName: aws.String(mfeTableName),
	}
	_, err = dynamoClient().PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return existing, nil
}

// Delete is a function
func Delete(momentID string) (*models.Moment, error) {
	existing, err := Read(momentID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"mfe_key": {
				S: aws.String(existing.MFEKey),
			},
		},
		TableName: aws.String(mfeTableName),
	}
	_, err = dynamoClient().DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return existing, nil
}
