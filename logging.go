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
	Entry    Entry    `json:"entry"`
	Messages []string `json:"messages"`
}

func (l *Log) Info(message string) {
	m := strings.Replace(message, "\n", "", -1)
	l.Messages = append(l.Messages, m)
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
	b, err := json.Marshal(l.Messages)
	if err == nil {
		l.Entry.Message = string(b)
	} else {
		return nil, err
	}

	b, err = json.Marshal(l.Entry)
	if err == nil {
		l.Messages = nil
	}
	return b, err
}

func Start() Log {
	now := time.Now()
	return Log{
		Entry: Entry{
			Timestamp: Timestamp{
				Seconds: now.Unix(),
				Nanos:   now.Nanosecond(),
			},
		},
	}
}
