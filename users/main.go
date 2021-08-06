package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tableName = "Users"

type User struct {
	UserId    string
	DeviceId  string
	FirstName string
	LastName  string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	userId := request.QueryStringParameters["userId"]
	deviceId := request.QueryStringParameters["deviceId"]
	fmt.Println("User: " + userId)
	fmt.Println("Device: " + deviceId)

	if userId == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 422,
		}, nil
	}

	user := getuser(userId, deviceId)

	if user == (User{}) {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, nil
	}

	formatedUser, err := json.Marshal(user)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(formatedUser),
		StatusCode: 200,
	}, nil
}

func getuser(userId string, deviceId string) User {

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials    <- make sure you declaring the profile when running locally or using default
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	var result *dynamodb.GetItemOutput
	var err error

	if deviceId == "" {
		result, err = svc.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"UserId": {
					S: aws.String(userId),
				},
				"DeviceId": {
					S: aws.String("1"),
				},
			},
		})
	} else {
		result, err = svc.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"UserId": {
					S: aws.String(userId),
				},
				"DeviceId": {
					S: aws.String(deviceId),
				},
			},
		})
	}

	if err != nil {
		log.Fatalf("Got error calling GetUser: %s", err)
	}

	user := User{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return user
}

func main() {
	lambda.Start(handler)
}
