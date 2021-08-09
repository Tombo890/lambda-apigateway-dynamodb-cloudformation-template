package main

import (
	"GO_fun/models"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockedGetItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.GetItemOutput
}

func (mockedOutput mockedGetItem) GetItem(in *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &mockedOutput.Response, nil
}

func TestHandler(t *testing.T) {
	userId := "1"
	deviceId := "1"
	firstName := "Some"
	lastName := "Guy"

	t.Run("Successful Request", func(t *testing.T) {
		mock := mockedGetItem{
			Response: dynamodb.GetItemOutput{
				Item: map[string]*dynamodb.AttributeValue{
					"UserId": {
						S: aws.String(userId),
					},
					"DeviceId": {
						S: aws.String(deviceId),
					},
					"FirstName": {
						S: aws.String(firstName),
					},
					"LastName": {
						S: aws.String(lastName),
					},
				},
			},
		}

		depend := dependencies{
			ddb:   mock,
			table: "test_table",
		}

		returnUser := depend.GetUser(userId, deviceId)

		if returnUser == (models.User{}) {
			t.Fatal("Something Wrong, panic!!!")
		} else if returnUser.DeviceId != deviceId {
			t.Fatal("DeviceId didn't match")
		} else if returnUser.FirstName != firstName {
			t.Fatal("FirstName didn't match")
		} else if returnUser.LastName != lastName {
			t.Fatal("LastName didn't match")
		} else if returnUser.UserId != userId {
			t.Fatal("UserId didn't match")
		}
	})
}
