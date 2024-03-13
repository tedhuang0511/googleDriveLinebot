package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func NewTableBasics(tableName string) *TableBasics {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	client := dynamodb.NewFromConfig(cfg)
	return &TableBasics{
		DynamoDbClient: client,
		TableName:      tableName,
	}
}

func CreateLocalClient(port int) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf("http://localhost:%d/", port)
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(dsn)
	})
}
