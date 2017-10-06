package main

import (
	"testing"

	"github.com/Cistern/cistern/internal/query"
)

func TestLimit(t *testing.T) {
	ec, err := CreateEventCollection("/tmp/test_cistern_limit.lm2")
	defer ec.col.Destroy()
	if err != nil {
		t.Fatal(err)
	}
	err = ec.StoreEvents(testEvents)
	if err != nil {
		t.Fatal(err)
	}

	const limit = 5
	result, err := ec.Query(query.Desc{
		Limit: limit,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Events) > limit {
		t.Errorf("expected at most %d events but got %d", limit, len(result.Events))
	}
}

func TestFilter(t *testing.T) {
	ec, err := CreateEventCollection("/tmp/test_cistern_limit.lm2")
	defer ec.col.Destroy()
	if err != nil {
		t.Fatal(err)
	}
	err = ec.StoreEvents(testEvents)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ec.Query(query.Desc{
		Filters: []query.Filter{
			{
				Column:    "source_port",
				Condition: "=",
				Value:     443.0,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Events) != 3 {
		t.Errorf("expected %d events but got %d", 3, len(result.Events))
	}
}
