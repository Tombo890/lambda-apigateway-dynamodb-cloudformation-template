package main

import (
	"GO_fun/models"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockedPutItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.PutItemOutput
}

func (mockedOutput mockedPutItem) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &mockedOutput.Response, nil
}

func TestHandler(t *testing.T) {

	deviceId := "1"
	firstName := "Some"
	lastName := "Guy"

	userWithoutIdToCreate := models.User{
		DeviceId:  deviceId,
		FirstName: firstName,
		LastName:  lastName,
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
		} else if returnUser.DeviceId != deviceId {
			t.Fatal("DeviceId didn't match")
		} else if returnUser.FirstName != firstName {
			t.Fatal("FirstName didn't match")
		} else if returnUser.LastName != lastName {
			t.Fatal("LastName didn't match")
		} else if returnUser.UserId == "" {
			t.Fatal("UserId wasn't created")
		}
	})
}