package slog

import "testing"

func TestAddLogMessage(t *testing.T) {
	log := Start()
	if log.entry.Timestamp.Seconds == 0 {
		t.Fatalf("log.entry.Timestamp.Seconds is Zero")
	}

	if log.entry.Timestamp.Nanos == 0 {
		t.Fatalf("log.entry.Timestamp.Nanos is Zero")
	}

	messages := []string{"Hello Logging", "Hello Logging Again !"}
	for i, m := range messages {
		log.Info(m)
		if len(log.messages) != i+1 {
			t.Fatalf("unexpected log.messages.len. %d != %d", len(log.messages), i)
		}
		if log.messages[i] != m {
			t.Fatalf("unexpected log.messages. %s != %s", log.messages[i], m)
		}
	}

	_, err := log.flush()
	if err != nil {
		t.Fatalf("log.flush() err. err = %s", err.Error())
	}
}
