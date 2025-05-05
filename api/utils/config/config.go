package config

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"runtime/debug"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// NewConfig returns a new config from filename passed
func NewConfig[T any](fileName string) (*T, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var config T
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type SecretString string

const (
	SecretsManagerURLScheme = "secrets-manager://"
)

func (ss *SecretString) UnmarshalJSON(data []byte) error {
	var ssString string
	if err := json.Unmarshal(data, &ssString); err != nil {
		return err
	}

	if !strings.HasPrefix(ssString, SecretsManagerURLScheme) {
		*ss = SecretString(ssString)
		return nil
	}

	u, err := url.Parse(ssString)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: fix the hard coded region
	ssString, err = SecretsManagerGetString(ctx, u.Host, "us-west-2")
	if err != nil {
		return err
	}

	*ss = SecretString(ssString)
	return nil
}

// TODO: eventually, probably don't want a new instance of this for each string
func SecretsManagerGetString(ctx context.Context, secretName string, region string) (string, error) {
	var opts []func(*awsConfig.LoadOptions) error
	opts = append(opts, awsConfig.WithRegion(region))
	if os.Getenv("ENVIRONMENT") == "local" {
		opts = append(opts, awsConfig.WithSharedConfigProfile("zero"))
	}
	cfg, err := awsConfig.LoadDefaultConfig(
		ctx,
		opts...,
	)
	if err != nil {
		debug.PrintStack()
		return "", err
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		debug.PrintStack()
		return "", err
	}

	if result.SecretString != nil {
		return *result.SecretString, nil
	}
	return "", nil
}
