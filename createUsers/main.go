package main

import (
	"GO_fun/models"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	jsoniter "github.com/json-iterator/go"
)

type dependencies struct {
	ddb   dynamodbiface.DynamoDBAPI
	table string
}

func (depend *dependencies) CreateUser(userToSave models.User) models.User {

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

	marshaledInput, err := dynamodbattribute.MarshalMap(userToSave)
	if err != nil {
		log.Fatalf("Failed to marshal new user: %s", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(depend.table),
		Item:      marshaledInput,
	}

	_, err = depend.ddb.PutItem(input)

	if err != nil {
		log.Fatalf("Got error calling GetUser: %s", err)
	}

	return userToSave
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	passedUser := models.User{}

	err := jsoniter.Unmarshal([]byte(request.Body), &passedUser)
	if err != nil || passedUser == (models.User{}) || passedUser.Email == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 422,
		}, nil
	}

	depend := dependencies{}
	userRecord := depend.CreateUser(passedUser)

	if userRecord == (models.User{}) {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, nil
	}

	formatedUser, err := jsoniter.Marshal(passedUser)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(formatedUser),
		StatusCode: 201,
	}, nil
}

func main() {

	lambda.Start(Handler)
}
