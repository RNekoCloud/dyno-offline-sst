package main

import (
	"context"
	"encoding/json"
	"log"
	"main/config"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Response struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}

type Dyno struct {
	Client *dynamodb.Client
}

func (dyno *Dyno) ListTables() []string {
	var tableName []string

	tables, err := dyno.Client.ListTables(context.Background(), &dynamodb.ListTablesInput{})

	if err != nil {
		log.Fatal("Failed to listing tables:", err)
	}

	tableName = tables.TableNames
	return tableName
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize Dyno configuration
	dynoClient := config.SetupDynoConfig()
	dyno := Dyno{
		Client: dynoClient,
	}

	table := dyno.ListTables()

	res := Response{
		Message: "Successfully retrieve list tables",
		Status:  "success",
		Data:    table,
	}

	b, _ := json.MarshalIndent(res, "", " ")

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application",
		},
		Body: string(b),
	}, nil

}

func main() {
	lambda.Start(Handler)
}
