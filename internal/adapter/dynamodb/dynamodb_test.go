package dynamodb

import (
	"os"
	"testing"
)

var testTableBasics *TableBasics

func TestMain(m *testing.M) {
	// google_oauth
	testTableBasics = NewTableBasics("google-oauth")
	// change to local dynamodb
	testTableBasics.DynamoDbClient = CreateLocalClient(8000)
	os.Exit(m.Run())
}
