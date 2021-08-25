package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/tidwall/gjson"
)

var (
	hostedZoneID = "" // NOTE: add your hosted zone id
)

func main() {

	sess := session.Must(session.NewSession())
	awsConf := aws.NewConfig()
	// ARN of the IAM role to assume
	// Remove the next two lines if you are not working with assumed roles or add your role
	creds := stscreds.NewCredentials(sess, "")
	awsConf.WithCredentials(creds)
	svc := route53.New(sess, awsConf)

	rrs, err := ListAllRecordSets(svc, hostedZoneID)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(rrs)
	json := string(b)
	if err != nil {
		log.Fatal(err)
	}
	if !gjson.Valid(json) {
		log.Fatal("invalid json")
	}
	data := gjson.Parse(json)
	data.ForEach(func(key, value gjson.Result) bool {
		raw := value.Raw
		name := getKey(&raw, "Name")
		typ := getKey(&raw, "Type")
		aliasTarget := getKey(&raw, "AliasTarget.DNSName")
		// INFO: Only prints A type records
		if typ == "A" {
			fmt.Println(name, aliasTarget)
		}
		return true // keep iterating
	})

}

// ListAllRecordSets returns a ResourceRecordSet containing all records
func ListAllRecordSets(r53 *route53.Route53, id string) (rrsets []*route53.ResourceRecordSet, err error) {
	req := route53.ListResourceRecordSetsInput{
		HostedZoneId: &id,
	}
	// iterate over pagination
	for {
		var resp *route53.ListResourceRecordSetsOutput
		resp, err = r53.ListResourceRecordSets(&req)
		if err != nil {
			return
		} else {
			rrsets = append(rrsets, resp.ResourceRecordSets...)
			if *resp.IsTruncated {
				req.StartRecordName = resp.NextRecordName
				req.StartRecordType = resp.NextRecordType
				req.StartRecordIdentifier = resp.NextRecordIdentifier
			} else {
				break
			}
		}
	}
	return
}

func getKey(json *string, key string) string {
	return gjson.Get(*json, key).String()
}
