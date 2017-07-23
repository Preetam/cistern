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

	limit := int64(10000)
	groupName := "flowlogs"

	eventCollection, err := OpenEventCollection("/tmp/flowlogs.lm2")
	if err != nil {
		if err == ErrDoesNotExist {
			eventCollection, err = CreateEventCollection("/tmp/flowlogs.lm2")
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	now := (time.Now().Unix() - 86400) * 1000
	lastEventTs := now
	var token *string
	var prevMessageID *string
	interleaved := true
	for {
		log.Println("starting batch ===================================================")
		output, err := svc.FilterLogEvents(&cloudwatchlogs.FilterLogEventsInput{
			LogGroupName: &groupName,
			Limit:        &limit,
			NextToken:    token,
			StartTime:    &lastEventTs,
			Interleaved:  &interleaved,
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

		collectionEvents := []Event{}
		for _, e := range output.Events {
			rec := &FlowLogRecord{}
			err = rec.Parse(*e.Message)
			if err == nil {
				rec.Timestamp = rec.Start
				rec.Duration = rec.End.Sub(rec.Start).Seconds()
				event := rec.ToEvent()
				event["_tag"] = *e.LogStreamName
				collectionEvents = append(collectionEvents, event)
			} else {
				log.Println(err)
			}
		}
		err = eventCollection.StoreEvents(collectionEvents)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("stored", len(collectionEvents), "events")

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
