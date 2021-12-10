package main

import (
	"GO_fun/models"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

type mockedPutItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.PutItemOutput
}

func (mockedOutput mockedPutItem) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &mockedOutput.Response, nil
}

func TestHandler(t *testing.T) {

	billingId := uuid.New().String()
	email := "test@paxi.com"
	plan := "Standard"
	phone := "1234567890"

	userToCreate := models.User{
		Email:     email,
		Plan:      plan,
		BillingId: billingId,
		Phone:     phone,
	}

	t.Run("Successful Request", func(t *testing.T) {
		mock := mockedPutItem{
			Response: dynamodb.PutItemOutput{},
		}

		depend := dependencies{
			ddb:   mock,
			table: "test_table",
		}

		returnUser := depend.CreateUser(userToCreate)

		if returnUser == (models.User{}) {
			t.Fatal("Something Wrong, panic!!!")
		}
	})
}
