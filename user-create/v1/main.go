package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/tech-talks/gotesting/internal/apigateway"
	"github.com/tech-talks/gotesting/user-create/v1/model"
	"github.com/tech-talks/gotesting/user-create/v1/repository"
)

type Request struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Age         *int    `json:"age"`
}

type Handler struct {
	userRepository userRepository
}

func NewHandler(
	userRepository userRepository,
) *Handler {
	return &Handler{
		userRepository,
	}
}

type userRepository interface {
	SaveUser(user model.User) error
}

func (h *Handler) Handle() func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		requestBody := Request{}
		err := json.Unmarshal([]byte(req.Body), &requestBody)
		if err != nil {
			return apigateway.Error(http.StatusBadRequest, err), nil
		}

		// request validation
		if strings.TrimSpace(requestBody.Email) == "" {
			return apigateway.Error(http.StatusBadRequest, errors.New("email is required")), nil
		}
		if strings.TrimSpace(requestBody.Name) == "" {
			return apigateway.Error(http.StatusBadRequest, errors.New("name is required")), nil
		}
		if requestBody.Age == nil {
			return apigateway.Error(http.StatusBadRequest, errors.New("age is required")), nil
		}

		phoneNumber := ""
		if requestBody.PhoneNumber != nil {
			phoneNumber = *requestBody.PhoneNumber
		}

		err = h.userRepository.SaveUser(model.User{
			Email:       requestBody.Email,
			Name:        requestBody.Name,
			Age:         *requestBody.Age,
			PhoneNumber: phoneNumber,
		})
		if err != nil {
			return apigateway.Error(http.StatusInternalServerError, errors.Wrap(err, "main: Handler.handle userRepository.SaveUser error")), nil
		}

		return apigateway.Respond(http.StatusOK, ""), nil
	}
}

func main() {
	dynamodbTableUser := os.Getenv("DYNAMODB_TABLE_USER")
	if dynamodbTableUser == "" {
		panic("DYNAMODB_TABLE_USER cannot be empty")
	}

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	dynamodbClient := dynamodb.New(sess)

	userRepository := repository.NewUserRepository(
		dynamodbClient,
		dynamodbTableUser,
	)

	handler := NewHandler(
		userRepository,
	)

	lambda.Start(handler.Handle())
}
