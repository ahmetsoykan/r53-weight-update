package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	svc := route53.New(sess)

	output, _ := svc.ListHostedZones(nil)

	var ids []string
	for i := range output.HostedZones {
		hostedZoneIds := strings.Split(*output.HostedZones[i].Id, "/")
		ids = append(ids, hostedZoneIds[2])
	}
	listRecords(svc, ids)

}

func listRecords(svc *route53.Route53, ids []string) {
	//fmt.Println(ids)

	var wint int64
	wint = 0
	for _, zoneid := range ids {
		zone, _ := svc.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(zoneid),
		})

		for id, record := range zone.ResourceRecordSets {
			if record.Weight != nil {
				if aws.Int64Value(record.Weight) == wint {
					//fmt.Println(*record.Name)
					//fmt.Println(zone.ResourceRecordSets[id])
					//fmt.Println(*zone.ResourceRecordSets[id].ResourceRecords[0].Value)
					//fmt.Println(zoneid)
					//fmt.Println(*zone.ResourceRecordSets[id].SetIdentifier)
					//fmt.Println(*record.Name)
					//fmt.Println("----------------------")
					updateRecord(svc, *record.Name, *zone.ResourceRecordSets[id].ResourceRecords[0].Value, *zone.ResourceRecordSets[id].SetIdentifier, zoneid)

				}
			}
		}
	}
}

func updateRecord(svc *route53.Route53, name string, target string, setid string, zoneid string) {

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{ // Required
			Changes: []*route53.Change{ // Required
				{ // Required
					Action: aws.String("UPSERT"), // Required
					ResourceRecordSet: &route53.ResourceRecordSet{ // Required
						Name: aws.String(name),    // Required
						Type: aws.String("CNAME"), // Required
						ResourceRecords: []*route53.ResourceRecord{
							{ // Required
								Value: aws.String(target), // Required
							},
						},
						TTL:           aws.Int64(300),
						Weight:        aws.Int64(5),
						SetIdentifier: aws.String(setid),
					},
				},
			},
			Comment: aws.String("Updated!"),
		},
		HostedZoneId: aws.String(zoneid), // Required
	}
	resp, err := svc.ChangeResourceRecordSets(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Change Response:")
	fmt.Println(resp)
}
