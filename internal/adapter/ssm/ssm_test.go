package ssm

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"os"
	"testing"
	"time"
)

var testSSM *SSM

func TestMain(m *testing.M) {

	testSSM = NewSSM()

	os.Exit(m.Run())
}

type SSMGetParameterImpl struct{}

func (dt SSMGetParameterImpl) GetParameter(ctx context.Context,
	params *ssm.GetParameterInput,
	optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {

	var parameter *types.Parameter

	if *params.Name == "secret-name" {
		parameter = &types.Parameter{Value: aws.String("secret-value")}
	}
	if *params.Name == "token-name" {
		parameter = &types.Parameter{Value: aws.String("token-value")}
	}

	output := &ssm.GetParameterOutput{
		Parameter: parameter,
	}

	return output, nil
}

type Config struct {
	MockChannelSecret      string `json:"MockChannelSecret"`
	MockChannelAccessToken string `json:"MockChannelAccessToken"`
}

var configFileName = "config.json"

var globalConfig Config

func populateConfiguration(t *testing.T) error {
	content, err := os.ReadFile(configFileName)
	if err != nil {
		return err
	}

	text := string(content)

	err = json.Unmarshal([]byte(text), &globalConfig)
	if err != nil {
		return err
	}

	if globalConfig.MockChannelSecret == "" {
		msg := "You must supply a value for MockChannelSecret in " + configFileName
		return errors.New(msg)
	}
	if globalConfig.MockChannelAccessToken == "" {
		msg := "You must supply a value for MockChannelAccessToken in " + configFileName
		return errors.New(msg)
	}

	return nil
}

func TestFindParameter(t *testing.T) {
	thisTime := time.Now()
	nowString := thisTime.Format("2006-01-02 15:04:05 Monday")
	t.Log("Starting unit test at " + nowString)

	err := populateConfiguration(t)
	if err != nil {
		t.Fatal(err)
	}

	api := &SSMGetParameterImpl{}

	respSecret, err := testSSM.FindParameter(context.Background(), *api, globalConfig.MockChannelSecret)
	if err != nil {
		t.Log("Got an error ...:")
		t.Log(err)
		return
	}

	t.Log("MockChannelSecret value: " + respSecret)

	respToken, err := testSSM.FindParameter(context.Background(), *api, globalConfig.MockChannelAccessToken)
	if err != nil {
		t.Log("Got an error ...:")
		t.Log(err)
		return
	}

	t.Log("MockChannelAccessToken value: " + respToken)
}
