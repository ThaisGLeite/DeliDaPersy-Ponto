// Package query provides utilities for interacting with DynamoDB to insert and select punch records.
package query

import (
	"context"
	"deli-ponto/pkg/model"
	"deli-ponto/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/beevik/ntp"
)

// TableName is the DynamoDB table name.
const TableName = "PontoColaborador"

// InsertPunch adds a punch record to the DynamoDB table.
// dynamoClient: The DynamoDB client
// nome: The name associated with the punch record
// logs: Logging utilities
func InsertPunch(dynamoClient *dynamodb.Client, nome string, logs *utils.GoAppTools) {
	// Fetch the current time from the National Observatory
	currentTime, err := ntp.Time("a.st1.ntp.br")
	logs.Check(err)

	// Create a new punch record
	punchRecord := model.Punch{
		Nome: nome,
		Data: currentTime.Format("2006-01-02_15:04:05"),
	}

	// Serialize the punch record to a map
	item, err := attributevalue.MarshalMap(punchRecord)
	logs.Check(err)

	// Insert the serialized record into DynamoDB
	_, err = dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})
	logs.Check(err)
}

// SelectPunch fetches a punch record from DynamoDB based on the name.
// Returns the punch record.
func SelectPunch(nome string, dynamoClient *dynamodb.Client, app *utils.GoAppTools) model.Punch {
	// Create a query expression
	queryExpr := expression.Name("Nome").Equal(expression.Value(nome))
	projectionExpr := expression.NamesList(expression.Name("Nome"), expression.Name("Data"))
	expr, err := expression.NewBuilder().WithFilter(queryExpr).WithProjection(projectionExpr).Build()
	app.Check(err) // Using Check method from utils

	// Create scan input parameters
	scanParams := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(TableName),
	}

	// Execute the query
	result, err := dynamoClient.Scan(context.TODO(), scanParams)
	app.CheckAndPanic(err)

	// Deserialize the results
	var punch model.Punch
	for _, item := range result.Items {
		err = attributevalue.UnmarshalMap(item, &punch)
		app.CheckAndPanic(err)
	}

	return punch
}

func SelectReport(nome string, periodo string, dynamoClient dynamodb.Client, logs *utils.GoAppTools) []model.Punch {
	ctx := context.Background()

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("PontoColaborador"),
		KeyConditionExpression: aws.String("#Nome = :nome AND begins_with(#Data, :mes)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nome": &types.AttributeValueMemberS{Value: nome},
			":mes":  &types.AttributeValueMemberS{Value: periodo},
		},
		ExpressionAttributeNames: map[string]string{
			"#Nome": "Nome",
			"#Data": "Data",
		},
	}
	queryOutput, err := dynamoClient.Query(ctx, queryInput)
	logs.Check(err)

	punchs := make([]model.Punch, 0)
	for _, item := range queryOutput.Items {
		punch := model.Punch{}
		err := attributevalue.UnmarshalMap(item, &punch)
		if err != nil {
			// Handle the error
			logs.Check(err)
			continue
		}
		punchs = append(punchs, punch)
	}

	return punchs
}
