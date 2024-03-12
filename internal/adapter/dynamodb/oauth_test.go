package dynamodb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateGoogleOAuthTable(t *testing.T) {
	tableDesc, err := testTableBasics.CreateGoogleOAuthTable()
	t.Log("tableDesc:", tableDesc)
	t.Log("ERROR:", err)
	assert.NoError(t, err, "Expected no error creating table")
	assert.NotNil(t, tableDesc, "Table description should not be nil")
}

func TestAddGoogleOAuthToken(t *testing.T) {
	tok := GoogleOAuthToken{
		PK:           "test12345",
		AccessToken:  "test123",
		TokenType:    "Bearer",
		RefreshToken: "test123",
		Expiry:       "2023-09-24T11:31:54.2936004+08:00",
	}
	err := testTableBasics.AddGoogleOAuthToken(tok)
	if err != nil {
		t.Log("ERROR:", err)
	}
}

func TestTxUpdateGoogleOAuthToken(t *testing.T) {
	tok := GoogleOAuthToken{
		PK:           "test1234",
		AccessToken:  "test77788",
		RefreshToken: "test77788",
	}
	output, err := testTableBasics.TxUpdateGoogleOAuthToken(tok)
	t.Log("output:", output)
	if err != nil {
		t.Log("ERROR:", err)
	}
}

func TestGetGoogleOAuthToken(t *testing.T) {
	tok, err := testTableBasics.GetGoogleOAuthToken("test1234")
	if err != nil {
		t.Log(err)
	}
	t.Log("Get Token:", tok)
}
