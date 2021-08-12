package repository

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/require"
	"github.com/tech-talks/gotesting/user-create/v1/model"
)

func Test_UserRepository_SaveUser(t *testing.T) {
	tableName := "user-table"
	closer, dynamodbLocal := DynamodbStart(t)
	defer closer()
	createUserTable(t, dynamodbLocal, tableName)

	repo := NewUserRepository(dynamodbLocal, tableName)

	err := repo.SaveUser(model.User{
		Email:       "email@server.com",
		Name:        "John Doe",
		PhoneNumber: "123123123",
		Age:         20,
	})
	require.NoError(t, err)

	foundUser, err := repo.FindUser("email@server.com")
	require.NoError(t, err)

	require.Equal(t, "John Doe", foundUser.Name)
	require.Equal(t, "123123123", foundUser.PhoneNumber)
	require.Equal(t, 20, foundUser.Age)
}

func createUserTable(t *testing.T, client *dynamodb.DynamoDB, tableName string) {
	_, err := client.CreateTable(&dynamodb.CreateTableInput{
		TableName: &tableName,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("email"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("email"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(2),
			WriteCapacityUnits: aws.Int64(2),
		},
	})
	require.NoError(t, err)
}
