package main

import (
	"testing"

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
	t.Run("Successful Request", func(t *testing.T) {
		mock := mockedGetItem{
			Response: dynamodb.GetItemOutput{
				//Item = map[string]
			},
		}

		depend := dependencies{
			ddb:   mock,
			table: "test_table",
		}

		returnUser := depend.GetUser("1", "1")
		if returnUser == (user{}) {
			t.Fatal("Something Wrong, panic!!!")
		}
	})
}
