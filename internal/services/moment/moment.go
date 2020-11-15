package moment

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	models "github.com/MomentsFromEarth/api/internal/models/moment"
	dynamodbmfe "github.com/MomentsFromEarth/api/internal/services/dynamodb"
	sqsmfe "github.com/MomentsFromEarth/api/internal/services/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/teris-io/shortid"
)

var mfeTableName = "MFE"
var mfeQuery01 = "query_key_01"
var mfeQuery02 = "query_key_02"
var momentJobQueueURL = "https://sqs.us-east-1.amazonaws.com/776913033148/moments.fifo"

func sqsClient() *sqs.SQS {
	return sqsmfe.Client()
}

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
		Captured:    newMoment.Captured,
		Created:     now,
		Updated:     now,
		QueryKey01:  getStatusQueryKey(status),
		QueryKey02:  "mom",
	}
}

func getMfeKey(momentID string) string {
	return fmt.Sprintf("mom:%s", momentID)
}

func getStatusQueryKey(status string) string {
	return fmt.Sprintf("mom:status:%s", strings.ToLower(status))
}

func momentJobParams(moment *models.Moment) *sqs.SendMessageInput {
	momJSON, _ := json.Marshal(moment)
	return &sqs.SendMessageInput{
		MessageBody:            aws.String(string(momJSON)),
		QueueUrl:               aws.String(momentJobQueueURL),
		MessageGroupId:         aws.String("mfe-api"),
		MessageDeduplicationId: aws.String(moment.QueueID),
	}
}

func addMomentToProcessingQueue(moment *models.Moment) {
	client := sqsClient()
	output, err := client.SendMessage(momentJobParams(moment))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Moment[%s] added to Processing Queue: %s", moment.MomentID, *output.MessageId)
	}
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
	addMomentToProcessingQueue(moment)
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
