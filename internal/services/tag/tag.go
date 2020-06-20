package tag

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	models "github.com/MomentsFromEarth/api/internal/models/tag"
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

func genTagID() (string, error) {
	return shortid.Generate()
}

func formatNewTag(name string) *models.Tag {
	now := genTimestamp()
	tagID, _ := genTagID()
	return &models.Tag{
		MFEKey:     getMfeKey(tagID),
		TagID:      tagID,
		Name:       name,
		Count:      0,
		Created:    now,
		Updated:    now,
		QueryKey01: getNameQueryKey(name),
		QueryKey02: "tag",
	}
}

func formatNameQueryInput(name string) *dynamodb.QueryInput {
	return &dynamodb.QueryInput{
		TableName:              aws.String(mfeTableName),
		IndexName:              aws.String(fmt.Sprintf("%s-index", mfeQuery01)),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :qk01", mfeQuery01)),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":qk01": {
				S: aws.String(getNameQueryKey(name)),
			},
		},
	}
}

func getMfeKey(tagID string) string {
	return fmt.Sprintf("tag:%s", tagID)
}

func getNameQueryKey(name string) string {
	return fmt.Sprintf("tag:name:%s", name)
}

// FilterAlphanumeric is a function
func filterAlphanumeric(str string) string {
	r := regexp.MustCompile("[^a-zA-Z0-9_-]+")
	return r.ReplaceAllString(str, "")
}

// FilterWhitespace is a function
func filterWhitespace(str string) string {
	return strings.Replace(str, " ", "", -1)
}

// FilterName is a function
func filterName(name string) string {
	return strings.ToLower(filterWhitespace(filterAlphanumeric(name)))
}

// Create is a function
func Create(newTag *models.NewTag) (*models.Tag, error) {
	name := filterName(newTag.Name)
	existing, _ := read(name)
	if existing != nil {
		existing.Count = existing.Count + 1
		tagInput, err := dynamodbattribute.MarshalMap(existing)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		input := &dynamodb.PutItemInput{
			Item:      tagInput,
			TableName: aws.String(mfeTableName),
		}
		_, err = dynamoClient().PutItem(input)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		return existing, nil
	} else {
		tag := formatNewTag(name)
		tagInput, err := dynamodbattribute.MarshalMap(tag)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		input := &dynamodb.PutItemInput{
			Item:      tagInput,
			TableName: aws.String(mfeTableName),
		}
		_, err = dynamoClient().PutItem(input)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		return tag, nil
	}
}

// Read is a function
func Read(tagID string) (*models.Tag, error) {
	dydb := dynamoClient()
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"mfe_key": {
				S: aws.String(getMfeKey(tagID)),
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
	tag := &models.Tag{}
	err = dynamodbattribute.UnmarshalMap(res, &tag)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return tag, nil
}

func read(name string) (*models.Tag, error) {
	n := filterName(name)
	dydb := dynamoClient()
	result, err := dydb.Query(formatNameQueryInput(n))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if len(result.Items) == 0 {
		msg := fmt.Errorf("Tag not found: %s", n)
		return nil, msg
	}

	res := result.Items[0]
	tag := &models.Tag{}
	err = dynamodbattribute.UnmarshalMap(res, &tag)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return tag, nil
}

// Delete is a function
func Delete(tagID string) (*models.Tag, error) {
	tag, err := Read(tagID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"mfe_key": {
				S: aws.String(getMfeKey(tagID)),
			},
		},
		TableName: aws.String(mfeTableName),
	}
	_, err = dynamoClient().DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return tag, nil
}
