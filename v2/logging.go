package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sinmetal/silverdog/dogtime"
)

const contextKey = "SINMETAL_SLOG"

// LogEntry is Stackdriver Logging Entry
type LogEntry struct {
	Timestamp   Timestamp `json:"timestamp"`
	Message     string    `json:"message"`
	Severity    string    `json:"severity"`
	severity    Severity
	Thread      int64       `json:"thread"`
	HttpRequest HttpRequest `json:"httpRequest"`
}

// Timestamp is Stackdriver Logging Timestamp
type Timestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int   `json:"nanos"`
}

type HttpRequest struct {
	RequestURL string `json:"requestUrl"`
}

// LogContainer is Log Object
type LogContainer struct {
	Entry    LogEntry `json:"entry"`
	Messages []string `json:"messages"`
}

type KV struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// WithVlaue is WithValue
func WithValue(ctx context.Context) context.Context {
	now := dogtime.Now()
	lc := createLogContainer(now)
	return context.WithValue(ctx, contextKey, lc)
}

// WithValueForHTTP is WithValue for HTTP Request
func WithValueForHTTP(ctx context.Context, r http.Request) context.Context {
	now := dogtime.Now()
	lc := createLogContainer(now)
	lc.Entry.HttpRequest = HttpRequest{
		RequestURL: r.URL.RawPath,
	}
	return context.WithValue(ctx, contextKey, lc)
}

// Value is return to LogEntry
func Value(ctx context.Context) (*LogContainer, bool) {
	v := ctx.Value(contextKey)
	if v == nil {
		return nil, false
	}
	l, ok := v.(*LogContainer)
	if !ok {
		return nil, false
	}
	return l, true
}

func Flush(ctx context.Context) {
	l, ok := Value(ctx)
	if !ok {
		log.Print("failed Flush Logging.context not include LogContainer\n")
		return
	}

	{
		j, err := json.Marshal(l.Messages)
		if err != nil {
			log.Printf("failed LogContainer.Messages to json.err=%+v,message=%+v\n", err, l.Messages)
		}
		l.Entry.Message = string(j)
	}

	j, err := json.Marshal(l.Entry)
	if err != nil {
		log.Printf("failed LogContainer to json.err=%+v,LogContainer=%+v\n", err, l)
	}
	fmt.Println(string(j))
}

func Info(ctx context.Context, message interface{}) {
	j, err := json.Marshal(message)
	if err != nil {
		log.Printf("failed log message to json.err=%+v,message=%+v\n", err, message)
	}

	l, ok := Value(ctx)
	if !ok {
		log.Printf("failed Info Logging.context not include LogContainer.message=%+v\n", message)
		return
	}
	l.Messages = append(l.Messages, string(j))
}

func createLogContainer(now time.Time) *LogContainer {
	return &LogContainer{
		Entry: LogEntry{
			Timestamp: Timestamp{
				Seconds: now.Unix(),
				Nanos:   now.Nanosecond(),
			},
			Thread: now.UnixNano(),
		},
	}
}
