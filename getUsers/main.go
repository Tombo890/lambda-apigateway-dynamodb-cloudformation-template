package main

import (
	"GO_fun/models"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
)

type dependencies struct {
	ddb   dynamodbiface.DynamoDBAPI
	table string
}

func (depend *dependencies) GetUser(email string, log *zap.Logger) models.User {

	if depend.ddb == nil {
		// Initialize a session that the SDK will use to load
		// credentials from the shared credentials file ~/.aws/credentials    <- make sure you declaring the profile when running locally or using default
		// and region from the shared configuration file ~/.aws/config.
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		// Create DynamoDB client
		svc := dynamodb.New(sess)

		depend = &dependencies{
			ddb:   svc,
			table: os.Getenv("TABLE"),
		}
	}

	result, err := depend.ddb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(depend.table),
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
	})

	if err != nil {
		log.Error("GetUser",
			zap.String("error message", err.Error()))
	}

	userRecord := models.User{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &userRecord)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return userRecord
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	email := request.QueryStringParameters["email"]
	logger.Info("GetUser",
		zap.String("email: ", email))

	if email == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 422,
		}, nil
	}

	depend := dependencies{}
	userRecord := depend.GetUser(email, logger)

	if userRecord == (models.User{}) {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, nil
	}

	formatedUser, err := jsoniter.Marshal(userRecord)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(formatedUser),
		StatusCode: 200,
	}, nil
}

func main() {

	lambda.Start(Handler)
}
