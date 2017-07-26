package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type FlowLogState struct {
	LastTimestamp int64  `json:"last_timestamp"`
	LastEventID   string `json:"last_event"`

	filename string
}

func (s *FlowLogState) Store() error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(s.filename, data, 0600)
	return err
}

func (s *FlowLogState) Load() error {
	data, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, s)
	return err
}

func NewFlowLogState(filename string) (*FlowLogState, error) {
	state := &FlowLogState{
		filename: filename,
	}
	state.Load()
	return state, state.Store()
}

// version account-id interface-id srcaddr dstaddr srcport dstport protocol packets bytes start end action log-status

type FlowLogRecord struct {
	Version       string    `json:"version"`
	AccountID     string    `json:"account_id"`
	InterfaceID   string    `json:"interface_id"`
	SourceAddress net.IP    `json:"source_address"`
	DestAddress   net.IP    `json:"dest_address"`
	SourcePort    int       `json:"source_port"`
	DestPort      int       `json:"dest_port"`
	Protocol      int       `json:"protocol"`
	Packets       int       `json:"packets"`
	Bytes         int       `json:"bytes"`
	Start         time.Time `json:"start"`
	End           time.Time `json:"end"`
	Action        string    `json:"action"`
	LogStatus     string    `json:"log_status"`

	Timestamp  time.Time `json:"_ts"`
	Duration   float64   `json:"_duration"`
	StreamName string    `json:"stream_name"`
}

func (r *FlowLogRecord) Parse(s string) error {
	parts := strings.Split(s, " ")
	if len(parts) != 14 {
		fmt.Println(parts)
		return errors.New("invalid flow log record")
	}

	r.Version = parts[0]
	r.AccountID = parts[1]
	r.InterfaceID = parts[2]
	r.SourceAddress = net.ParseIP(parts[3])
	r.DestAddress = net.ParseIP(parts[4])

	n, err := strconv.ParseInt(parts[5], 10, 64)
	if err != nil {
		return err
	}

	r.SourcePort = int(n)

	n, err = strconv.ParseInt(parts[6], 10, 64)
	if err != nil {
		return err
	}
	r.DestPort = int(n)

	n, err = strconv.ParseInt(parts[7], 10, 64)
	if err != nil {
		return err
	}
	r.Protocol = int(n)

	n, err = strconv.ParseInt(parts[8], 10, 64)
	if err != nil {
		return err
	}
	r.Packets = int(n)

	n, err = strconv.ParseInt(parts[9], 10, 64)
	if err != nil {
		return err
	}
	r.Bytes = int(n)

	n, err = strconv.ParseInt(parts[10], 10, 64)
	if err != nil {
		return err
	}
	r.Start = time.Unix(n, 0).UTC()

	n, err = strconv.ParseInt(parts[11], 10, 64)
	if err != nil {
		return err
	}
	r.End = time.Unix(n, 0).UTC()

	r.Action = parts[12]
	r.LogStatus = parts[13]

	return nil
}

func (r *FlowLogRecord) ToEvent() Event {
	return Event{
		"version":        r.Version,
		"account_id":     r.AccountID,
		"interface_id":   r.InterfaceID,
		"source_address": r.SourceAddress,
		"dest_address":   r.DestAddress,
		"source_port":    r.SourcePort,
		"dest_port":      r.DestPort,
		"protocol":       r.Protocol,
		"packets":        r.Packets,
		"bytes":          r.Bytes,
		"start":          r.Start,
		"end":            r.End,
		"action":         r.Action,
		"log_status":     r.LogStatus,
		"_ts":            r.Timestamp.Format(time.RFC3339Nano),
		"_duration":      r.Duration,
		"_tag":           r.StreamName,
	}
}

func groupFlowRecords(records []*FlowLogRecord) []FlowLogRecord {
	type groupKey struct {
		Timestamp     time.Time
		SourceAddress [16]byte
		DestAddress   [16]byte
		SourcePort    int
		DestPort      int
		Protocol      int
	}
	groups := map[groupKey]*FlowLogRecord{}
	for _, rec := range records {
		key := groupKey{
			Timestamp:     rec.Timestamp.Truncate(time.Minute * 10),
			SourceAddress: ipTo16Bytes(rec.SourceAddress),
			DestAddress:   ipTo16Bytes(rec.DestAddress),
			SourcePort:    rec.SourcePort,
			DestPort:      rec.DestPort,
			Protocol:      rec.Protocol,
		}
		groupRec := groups[key]
		if groupRec == nil {
			groupRec = &FlowLogRecord{
				Timestamp:     key.Timestamp,
				SourceAddress: net.IP(key.SourceAddress[:]),
				DestAddress:   net.IP(key.DestAddress[:]),
				SourcePort:    key.SourcePort,
				DestPort:      key.DestPort,
				Protocol:      key.Protocol,
			}
			groups[key] = groupRec
		}
		groupRec.Bytes += rec.Bytes
		groupRec.Packets += rec.Packets
	}

	result := []FlowLogRecord{}
	for _, rec := range groups {
		result = append(result, *rec)
	}
	return result
}

func ipTo16Bytes(ip net.IP) [16]byte {
	result := [16]byte{}
	copy(result[:], ip.To16())
	return result
}

func captureFlowLogs(groupName string, done chan struct{}) error {
	stop := make(chan struct{}, 1)

	go func() {
		<-done
		stop <- struct{}{}
	}()

	sess := session.Must(session.NewSession())
	svc := cloudwatchlogs.New(sess)

	limit := int64(10000)

	stateFile, err := NewFlowLogState(filepath.Join(DataDir, groupName+".state"))
	if err != nil {
		return err
	}

	collectionsLock.Lock()
	eventCollection := Collections[groupName]
	if eventCollection == nil {
		eventCollection, err = OpenEventCollection(filepath.Join(DataDir, groupName+".lm2"))
		if err != nil {
			if err == ErrDoesNotExist {
				eventCollection, err = CreateEventCollection(filepath.Join(DataDir, groupName+".lm2"))
			}
			if err != nil {
				collectionsLock.Unlock()
				return err
			}
		}
		Collections[groupName] = eventCollection
	}
	collectionsLock.Unlock()

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
	timer := time.NewTimer(0)
	for {

		select {
		case <-timer.C:
			log.Println("tick")
		case <-stop:
			log.Println("stopping", groupName)
			eventCollection.col.Close()
			return nil
		}

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
			return err
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
				return err
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
				event["_hash"] = "000000"
				events = append(events, event)
			}
			err = eventCollection.StoreEvents(events)
			if err != nil {
				return err
			}

			if len(output.Events) > 0 {
				stateFile.LastTimestamp = nextBatchStart
				stateFile.Store()
			}
			next := time.Unix(batchEnd/1000+batchSizeSeconds, 0)
			if time.Now().After(next) {
				timer.Reset(0)
			} else {
				timer.Reset(next.Sub(time.Now()))
			}
			nextBatchStart = batchEnd + 1
			token = nil
			currentBatch = currentBatch[:0]
		} else {
			token = output.NextToken
			timer.Reset(0)
		}
	}
}

func sleepUntil(ts time.Time) {
	if time.Now().After(ts) {
		return
	}
	time.Sleep(ts.Sub(time.Now()))
}
