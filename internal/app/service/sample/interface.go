package sample

import "firstProject/internal/adapter/dynamodb"

type SampleServiceDynamodbI interface {
	GetGoogleOAuthToken(line_userid string) (dynamodb.GoogleOAuthToken, error)
}
