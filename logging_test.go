package slog

import (
	"strings"
	"testing"
)

func TestAddLogMessage(t *testing.T) {
	log := Start()
	if log.Entry.Timestamp.Seconds == 0 {
		t.Fatalf("log.entry.Timestamp.Seconds is Zero")
	}

	if log.Entry.Timestamp.Nanos == 0 {
		t.Fatalf("log.entry.Timestamp.Nanos is Zero")
	}

	messages := []string{"Good morning Logging", "Hello Logging", "Good evening"}
	for i, m := range messages {
		log.Info(m)
		if len(log.Messages) != i+1 {
			t.Fatalf("unexpected log.messages.len. %d != %d", len(log.Messages), i)
		}
		if log.Messages[i] != m {
			t.Fatalf("unexpected log.messages. %s != %s", log.Messages[i], m)
		}
	}

	b, err := log.flush()
	if err != nil {
		t.Fatalf("log.flush() err. err = %s", err.Error())
	}
	if strings.Contains(string(b), "Good morning") == false {
		t.Fatalf("output json not contains Good morning")
	}
	if strings.Contains(string(b), "Hello") == false {
		t.Fatalf("output json not contains Hello")
	}
	if strings.Contains(string(b), "Good evening") == false {
		t.Fatalf("output json not contains Good evening")
	}
}
