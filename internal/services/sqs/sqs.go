package sqsmfe

import (
	awsmfe "github.com/MomentsFromEarth/api/internal/services/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var client *sqs.SQS

// Client returns AWS SQS client
func Client() *sqs.SQS {
	if client == nil {
		session := awsmfe.Session()
		client = sqs.New(session)
	}
	return client
}
