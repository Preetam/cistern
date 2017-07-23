package main

import (
	"encoding/json"
	"fmt"
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

	now := (time.Now().Unix() - 6000) * 1000
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

		for _, e := range output.Events {
			rec := &FlowLogRecord{}
			err = rec.Parse(*e.Message)
			if err == nil {
				rec.Timestamp = rec.Start
				rec.Duration = rec.End.Sub(rec.Start).Seconds()
				b, _ := json.Marshal(rec)
				fmt.Println(string(b))
			} else {
				log.Println(err)
			}
		}

		if output.NextToken == nil {
			time.Sleep(60 * time.Second)
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
