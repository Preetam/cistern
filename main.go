package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	sess := session.Must(session.NewSession())
	svc := cloudwatchlogs.New(sess)

	limit := int64(10000)
	groupName := "flowlogs"

	stateFile, err := NewFlowLogState("/tmp/flow.state")
	if err != nil {
		log.Fatal(err)
	}

	eventCollection, err := OpenEventCollection("/tmp/flowlogs.lm2")
	if err != nil {
		if err == ErrDoesNotExist {
			eventCollection, err = CreateEventCollection("/tmp/flowlogs.lm2")
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println("Got signal", sig)
		log.Println("Closing collection")
		eventCollection.col.Close()
		os.Exit(0)
	}()

	const batchSizeSeconds = 3600
	const batchSizeMilliseconds = batchSizeSeconds * 1000

	now := stateFile.LastTimestamp
	if now == 0 {
		now = (time.Now().Unix() - 86400) * 1000
	}
	nextBatchStart := (now / batchSizeMilliseconds) * batchSizeMilliseconds
	var token *string
	interleaved := true

	currentBatch := []*FlowLogRecord{}
	for {
		log.Println("starting batch", nextBatchStart, "===================================================")
		batchEnd := nextBatchStart + batchSizeMilliseconds - 1
		output, err := svc.FilterLogEvents(&cloudwatchlogs.FilterLogEventsInput{
			LogGroupName: &groupName,
			Limit:        &limit,
			NextToken:    token,
			StartTime:    &nextBatchStart,
			EndTime:      &batchEnd,
			Interleaved:  &interleaved,
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range output.Events {
			rec := &FlowLogRecord{}
			err = rec.Parse(*e.Message)
			if err == nil {
				rec.Timestamp = rec.Start
				rec.Duration = rec.End.Sub(rec.Start).Seconds()
				rec.StreamName = *e.LogStreamName
				currentBatch = append(currentBatch, rec)
			} else {
				log.Println(err)
			}
		}

		if output.NextToken == nil {
			log.Println("raw events:", len(currentBatch))
			grouped := groupFlowRecords(currentBatch)
			log.Println("grouped events:", len(grouped))
			events := []Event{}
			for _, rec := range grouped {
				event := rec.ToEvent()
				event["_tag"] = "flowlog"
				events = append(events, event)
			}
			err = eventCollection.StoreEvents(events)
			if err != nil {
				log.Fatal(err)
			}

			if len(output.Events) > 0 {
				stateFile.LastTimestamp = nextBatchStart
				stateFile.Store()
			}
			sleepUntil(time.Unix(batchEnd/1000+batchSizeSeconds, 0))
			nextBatchStart = batchEnd + 1
			token = nil
			currentBatch = currentBatch[:0]
		} else {
			token = output.NextToken
		}
	}
}

func sleepUntil(ts time.Time) {
	if time.Now().After(ts) {
		return
	}
	time.Sleep(ts.Sub(time.Now()))
}
