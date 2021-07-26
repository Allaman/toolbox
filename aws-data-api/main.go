package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

func main() {
	sess := session.Must(session.NewSession())
	awsConf := aws.NewConfig()
	awsConf.WithRegion("eu-central-1")
	// ARN of the IAM role to assume
	// Remove the next two lines if you are not working with assumed roles
	creds := stscreds.NewCredentials(sess, "")
	awsConf.WithCredentials(creds)
	svc := rdsdataservice.New(sess, awsConf)

	// Name of the database of the Aurora serverless cluster to query to
	database := ""
	// ARN of the Aurora serverless cluster
	resourceArn := ""
	// ARN of the secrets in AWS secrets manager containing database credentials
	secretArn := ""
	// SQL to be executed
	sql := "show tables;"
	input := rdsdataservice.ExecuteStatementInput{
		Database:    &database,
		ResourceArn: &resourceArn,
		SecretArn:   &secretArn,
		Sql:         &sql,
	}
	out, err := svc.ExecuteStatement(&input)
	if err != nil {
		log.Fatal(err)
	}
	// Printing raw response
	fmt.Println(out)
}
