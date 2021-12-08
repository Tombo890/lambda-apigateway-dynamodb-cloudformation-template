package main

import (
	"GO_fun/models"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

type mockedGetItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.GetItemOutput
}

func (mockedOutput mockedGetItem) GetItem(in *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &mockedOutput.Response, nil
}

func TestHandler(t *testing.T) {
	userId := uuid.New().String()
	email := "test@paxi.com"
	plan := "Standard"
	billingId := uuid.New().String()

	t.Run("Successful Request", func(t *testing.T) {
		mock := mockedGetItem{
			Response: dynamodb.GetItemOutput{
				Item: map[string]*dynamodb.AttributeValue{
					"UserId": {
						S: aws.String(userId),
					},
					"Email": {
						S: aws.String(email),
					},
					"Plan": {
						S: aws.String(plan),
					},
					"BillingId": {
						S: aws.String(billingId),
					},
				},
			},
		}

		depend := dependencies{
			ddb:   mock,
			table: "test_table",
		}

		returnUser := depend.GetUser(userId)

		if returnUser == (models.User{}) {
			t.Fatal("Something Wrong, panic!!!")
		} else if returnUser.Id != userId {
			t.Fatal("UserId didn't match")
		}
	})
}
