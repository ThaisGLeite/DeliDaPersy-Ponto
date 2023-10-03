// Package driver provides functionalities to set up the AWS configuration and DynamoDB client.
package driver

import (
	"context"
	"fmt"

	"deli-ponto/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// ConfigAws initializes and returns a new DynamoDB client with AWS configuration.
// It uses the shared credentials and config files located in the 'driver/data/' directory.
// Returns:
// - A pointer to the DynamoDB client
// - An error, if any occurred during the initialization
func ConfigAws(logs *utils.GoAppTools) (*dynamodb.Client, error) {
	// Create a context
	ctx := context.TODO()

	// Load AWS configuration from shared credentials and config files
	configAws, err := config.LoadDefaultConfig(
		ctx,
		config.WithSharedCredentialsFiles([]string{"pkg/driver/data/credentials.aws"}),
		config.WithSharedConfigFiles([]string{"pkg/driver/data/config.aws"}),
	)
	if err != nil {
		logs.Check(fmt.Errorf("failed to load AWS configuration: %w", err))
		return nil, err
	}

	// Initialize a new DynamoDB client from the AWS configuration
	dynamoClient := dynamodb.NewFromConfig(configAws)

	return dynamoClient, nil
}
