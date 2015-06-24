package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)

type Route53Zone struct {
	Name string
	Id string
	Records []Route53Record
}

type Route53Record struct {
	Name string
	Weight int64
	Type string
	AliasTarget string
	SetId string
	Ttl int64
}

func recordFactory(record *ListResourceRecordSetOutput) Route53Record{
	/*
	Creates a Route53Record from the aws response struct
	 */
	var tmpRecord Route53Record
	return tmpRecord
}

func recordsByZone(name *string) *[]Route53Record {
	var records []Route53Record
	return &records
}

func zoneByName(name *string) *Route53Zone {
	var zone Route53Zone
	svc := route53.New(nil)

	params := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(*name),
	}

	resp, err := svc.ListHostedZonesByName(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(awsutil.StringValue(resp))

	if resp.HostedZones != nil {
		/*
		There is actually a record.  i.e. this is a real domain
		 */
		zone.Name = *resp.DNSName

		/*
		Get the id from the hosted zone info
		 */
		zone.Id = *resp.HostedZones[0].ID
	}

	return &zone
}

func main() {
	zones := make(map[string]Route53Zone)
	svc := route53.New(&aws.Config{Region: "us-east-1"})

	resp, err := svc.ListHostedZones(nil)
	if err != nil {
		panic(err)
	}
//	fmt.Println(awsutil.StringValue(resp))

	for _, zone := range resp.HostedZones {
		var tmpRecords []Route53Record
		var zoneInfo Route53Zone
		zoneInfo.Name = *zone.Name
		zoneInfo.Id = *zone.ID

		params := &route53.ListResourceRecordSetsInput{
			HostedZoneID:	aws.String(*zone.ID),
		}

		resp, err := svc.ListResourceRecordSets(params)
		if err != nil {
			panic(err)
		}

//		fmt.Println(awsutil.StringValue(resp))
		for _, record := range resp.ResourceRecordSets {
			var tmpRecord Route53Record
			tmpRecord.Name = *record.Name
			tmpRecord.Type = *record.Type

			if record.TTL != nil {
				tmpRecord.Ttl = *record.TTL
			}
			if record.SetIdentifier != nil {
				tmpRecord.SetId = *record.SetIdentifier
			}
			if record.Weight != nil {
				tmpRecord.Weight = *record.Weight
			} else {
				tmpRecord.Weight = 0
			}
			if record.AliasTarget != nil {
				tmpRecord.AliasTarget = *record.AliasTarget.DNSName
			}
			tmpRecords = append(tmpRecords, tmpRecord)
		}
		zoneInfo.Records = tmpRecords
		zones[*zone.Name] = zoneInfo
	}
//	fmt.Println(zones)
	name := "mode.st"
	temp := zoneByName(&name)
	if temp.Name != name {
		panic("ahhh")
	}
}