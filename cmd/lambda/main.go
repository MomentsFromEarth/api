package main

import (
	"context"
	"log"

	api "github.com/MomentsFromEarth/api/internal"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Printf("lambda.init")
	apiEngine := api.Init()
	ginLambda = ginadapter.New(apiEngine)
}

// Handler is the Lambda entrypoint
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
