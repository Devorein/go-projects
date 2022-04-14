package handlers

import (
	"net/http"

	"go-serverless/pkg/user"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var ErrorMethodNotAllowed = "method not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]

	if len(email) > 0 {
		result, err := user.GetUser(email, tableName, dynaClient)
		if err != nil {
			return ApiResponse(http.StatusBadRequest, ErrorBody{
				ErrorMsg: aws.String(err.Error()),
			})
		}
		return ApiResponse(http.StatusOK, result)
	} else {
		users, err := user.GetUsers(tableName, dynaClient)

		if err != nil {
			return ApiResponse(http.StatusBadRequest, ErrorBody{
				ErrorMsg: aws.String(err.Error()),
			})
		}

		return ApiResponse(http.StatusOK, users)
	}
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	user, err := user.CreateUser(req, tableName, dynaClient)

	if err != nil {
		return ApiResponse(http.StatusBadRequest, ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}

	return ApiResponse(http.StatusOK, user)
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	user, err := user.UpdateUser(req, tableName, dynaClient)

	if err != nil {
		return ApiResponse(http.StatusBadRequest, ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}

	return ApiResponse(http.StatusOK, user)
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	err := user.DeleteUser(req, tableName, dynaClient)

	if err != nil {
		return ApiResponse(http.StatusBadRequest, ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}

	return ApiResponse(http.StatusOK, nil)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return ApiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
