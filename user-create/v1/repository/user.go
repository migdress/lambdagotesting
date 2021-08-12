package repository

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/tech-talks/gotesting/user-create/v1/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	dynamodbClient *dynamodb.DynamoDB
	tableName      string
}

func NewUserRepository(
	dynamodbClient *dynamodb.DynamoDB,
	tableName string,
) *UserRepository {
	return &UserRepository{
		dynamodbClient,
		tableName,
	}
}

func (r *UserRepository) SaveUser(user model.User) error {
	_, err := r.dynamodbClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(user.Email),
			},
			"name": {
				S: aws.String(user.Name),
			},
			"phone_number": {
				S: aws.String(user.PhoneNumber),
			},
			"age": {
				N: aws.String(strconv.Itoa(user.Age)),
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "repository: UserRepository.SaveUser dynamodbClient.PutItem error")
	}

	return nil
}

func (r *UserRepository) FindUser(email string) (*model.User, error) {
	res, err := r.dynamodbClient.Query(&dynamodb.QueryInput{
		TableName: aws.String(r.tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"email": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(email),
					},
				},
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "repository: UserRepository.FindUser dynamodbClient.Query error")
	}

	if len(res.Items) == 0 {
		return nil, errors.Wrap(ErrUserNotFound, "repository: UserRepository.FindUser error")
	}
	user := model.User{}
	err = dynamodbattribute.UnmarshalMap(res.Items[0], &user)

	return &user, nil
}
