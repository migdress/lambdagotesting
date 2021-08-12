package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"github.com/tech-talks/gotesting/user-create/v1/model"
)

func Test_Handler_Handle(t *testing.T) {
	testCases := []struct {
		name               string
		request            string
		userRepository     *mockUserRepository
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "happy path",
			request: `{
				"email": "blah@algo.com",
				"name": "blah",
				"age": 20
			}`,
			userRepository:     &mockUserRepository{},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewHandler(tc.userRepository)
			apiGatewayRequest := events.APIGatewayProxyRequest{
				Body: tc.request,
			}

			response, err := handler.Handle()(nil, apiGatewayRequest)
			require.NoError(t, err)

			require.Equal(t, tc.expectedStatusCode, response.StatusCode)
			require.Equal(t, tc.expectedResponse, response.Body)
		})
	}
}

type mockUserRepository struct {
	errSaveUser error
}

func (m *mockUserRepository) SaveUser(user model.User) error {
	return m.errSaveUser
}
