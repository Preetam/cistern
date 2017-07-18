package main

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func main() {
	sess := session.Must(session.NewSession())
	svc := cloudwatchlogs.New(sess)

	limit := int64(1000)
	groupName := "flowlogs"

	now := (time.Now().Unix() - 600) * 1000
	lastEventTs := now
	var token *string
	var prevMessageID *string
	for {
		log.Println("starting batch ===================================================")
		output, err := svc.FilterLogEvents(&cloudwatchlogs.FilterLogEventsInput{
			LogGroupName: &groupName,
			Limit:        &limit,
			NextToken:    token,
			StartTime:    &lastEventTs,
		})
		if err != nil {
			log.Fatal(err)
		}

		if prevMessageID != nil {
			for i, event := range output.Events {
				if *event.EventId == *prevMessageID {
					output.Events = output.Events[i+1:]
					break
				}
			}
		}

		log.Println(output.Events)

		if output.NextToken == nil {
			time.Sleep(5 * time.Second)
			token = nil
			if len(output.Events) > 0 {
				lastEventTs = *output.Events[len(output.Events)-1].Timestamp
				prevMessageID = output.Events[len(output.Events)-1].EventId
			}
		} else {
			token = output.NextToken
		}
	}
}
