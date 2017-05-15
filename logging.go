package slog

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Entry struct {
	Timestamp Timestamp `json:"timestamp"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	Thread    int       `json:"thread"`
}

type Timestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int   `json:"nanos"`
}

type Log struct {
	entry    Entry
	messages []string
}

func (l *Log) Info(message string) {
	m := strings.Replace(message, "\n", "", -1)
	l.messages = append(l.messages, m)
}

func (l *Log) Flush() {
	j, err := l.flush()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%s\n", j)
}

func (l *Log) flush() ([]byte, error) {
	b, err := json.Marshal(l.messages)
	if err == nil {
		l.entry.Message = string(b)
	} else {
		return nil, err
	}

	b, err = json.Marshal(l.entry)
	if err == nil {
		l.messages = nil
	}
	return b, err
}

func Start() Log {
	now := time.Now()
	return Log{
		entry: Entry{
			Timestamp: Timestamp{
				Seconds: now.Unix(),
				Nanos:   now.Nanosecond(),
			},
		},
	}
}
