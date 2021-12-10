package main

import (
	"GO_fun/models"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type mockedGetItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.GetItemOutput
}

func (mockedOutput mockedGetItem) GetItem(in *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &mockedOutput.Response, nil
}

func TestHandler(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	billingId := uuid.New().String()
	email := "test@paxi.com"
	plan := "Standard"
	phone := "1234567890"

	t.Run("Successful Request", func(t *testing.T) {
		mock := mockedGetItem{
			Response: dynamodb.GetItemOutput{
				Item: map[string]*dynamodb.AttributeValue{

					"BillingId": {
						S: aws.String(billingId),
					},
					"Email": {
						S: aws.String(email),
					},
					"Plan": {
						S: aws.String(plan),
					},
					"Phone": {
						S: aws.String(phone),
					},
				},
			},
		}

		depend := dependencies{
			ddb:   mock,
			table: "test_table",
		}

		returnUser := depend.GetUser(email, logger)

		if returnUser == (models.User{}) {
			t.Fatal("Something Wrong, panic!!!")
		} else if returnUser.Email != email {
			t.Fatal("UserId didn't match")
		}
	})
}
