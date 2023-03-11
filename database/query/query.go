package query

import (
	"context"
	"deli-ponto/configuration"
	"deli-ponto/model"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/beevik/ntp"
)

func InsertPunch(dynamoClient *dynamodb.Client, nome string, logs configuration.GoAppTools) {

	//o codigo esta indo no observatorio nacional pegar a data e hora
	datatemp, err := ntp.Time("a.st1.ntp.br")
	configuration.Check(err, logs)

	//formatando a data retornada do observatorio para a data no formato desejado (yy-mm-dd_hh:mm)
	tempY := datatemp.Format("06")
	tempM := datatemp.Format("01")
	tempD := datatemp.Format("02")
	tempH := strconv.Itoa(datatemp.Hour())
	tempMin := strconv.Itoa(datatemp.Minute())

	ponto := model.Punch{
		Nome: nome,
		Data: tempY + "-" + tempM + "-" + tempD + "_" + tempH + ":" + tempMin,
	}

	item, err := attributevalue.MarshalMap(ponto)

	configuration.Check(err, logs)

	_, err = dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("PontoColaborador"),
		Item:      item,
	})
	configuration.Check(err, logs)
}

func SelectPunch(Nome string, dynamoClient dynamodb.Client, app configuration.GoAppTools) model.Punch {
	query := expression.Name("Nome").Equal(expression.Value(Nome))
	proj := expression.NamesList(expression.Name("Nome"), expression.Name("Data"))

	expr, err := expression.NewBuilder().WithFilter(query).WithProjection(proj).Build()
	configuration.Check(err, app)

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("PontoColaborador"),
	}

	// Make the DynamoDB Query API call
	result, err := dynamoClient.Scan(context.TODO(), params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}

	var punch model.Punch
	for _, i := range result.Items {
		item := model.Punch{}

		err = attributevalue.UnmarshalMap(i, &item)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)

		}
		punch = item
	}

	return punch
}
