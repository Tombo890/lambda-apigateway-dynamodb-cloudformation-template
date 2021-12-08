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

	email := "test@paxi.com"
	plan := "Standard"
	billingId := uuid.New().String()

	userWithoutIdToCreate := models.User{
		Email:     email,
		Plan:      plan,
		BillingId: billingId,
	}

	t.Run("Successful Request", func(t *testing.T) {
		mock := mockedPutItem{
			Response: dynamodb.PutItemOutput{},
		}

		depend := dependencies{
			ddb:   mock,
			table: "test_table",
		}

		returnUser := depend.CreateUser(userWithoutIdToCreate)

		if returnUser == (models.User{}) {
			t.Fatal("Something Wrong, panic!!!")
		} else if returnUser.Id == "" {
			t.Fatal("UserId wasn't created")
		}
	})
}
