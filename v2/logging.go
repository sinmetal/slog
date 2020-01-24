package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sinmetal/gcpmetadata"
	"github.com/sinmetal/silverdog/dogtime"
	"go.opencensus.io/trace"
)

const contextKey = "SINMETAL_SLOG"

// LogEntry is Stackdriver Logging Entry
type LogEntry struct {
	LogName     string            `json:"logName"`
	Timestamp   Timestamp         `json:"timestamp"`
	TraceID     string            `json:"trace,omitempty"`
	SpanID      string            `json:"spanId,omitempty"`
	Operation   LogEntryOperation `json:"operation,omitempty"`
	Messages    []string          `json:"messages"`
	Severity    string            `json:"severity,omitempty"`
	severity    Severity
	HTTPRequest *HTTPRequest `json:"httpRequest,omitempty"`
}

// Timestamp is Stackdriver Logging Timestamp
type Timestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int   `json:"nanos"`
}

// HTTPRequest provides HTTPRequest log.
// spec: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#httprequest
type HTTPRequest struct {
	RequestMethod                  string `json:"requestMethod"`
	RequestURL                     string `json:"requestUrl"`
	RequestSize                    int64  `json:"requestSize,string,omitempty"`
	Status                         int    `json:"status,omitempty"`
	ResponseSize                   int64  `json:"responseSize,string,omitempty"`
	UserAgent                      string `json:"userAgent,omitempty"`
	RemoteIP                       string `json:"remoteIp,omitempty"`
	Referer                        string `json:"referer,omitempty"`
	Latency                        string `json:"latency,omitempty"`
	CacheLookup                    *bool  `json:"cacheLookup,omitempty"`
	CacheHit                       *bool  `json:"cacheHit,omitempty"`
	CacheValidatedWithOriginServer *bool  `json:"cacheValidatedWithOriginServer,omitempty"`
	CacheFillBytes                 *int64 `json:"cacheFillBytes,string,omitempty"`
	Protocol                       string `json:"protocol"`
}

// LogEntryOperation is Additional information about a potentially long-running operation with which a log entry is associated.
// spec: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logentryoperation
type LogEntryOperation struct {
	ID       string `json:"id"`
	Producer string `json:"producer"`
	First    bool   `json:"first"`
	Last     bool   `json:"last"`
}

// LogContainer is Log Object
type LogContainer struct {
	Entry LogEntry `json:"entry"`
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
func WithValueForHTTP(ctx context.Context, r *http.Request) context.Context {
	now := dogtime.Now()
	lc := createLogContainer(now)

	span := trace.FromContext(r.Context())
	lc.Entry.TraceID = span.SpanContext().TraceID.String()
	lc.Entry.HTTPRequest = &HTTPRequest{
		RequestMethod: r.Method,
		RequestURL:    r.RequestURI,
		UserAgent:     r.UserAgent(),
		RemoteIP:      r.RemoteAddr,
		Referer:       r.Referer(),
		Protocol:      r.Proto,
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

	j, err := json.Marshal(l.Entry)
	if err != nil {
		log.Printf("failed LogContainer to json.err=%+v,LogContainer=%+v\n", err, l)
	}
	fmt.Println(string(j))
}

// HTTPResponse is FlushWithHTTPResponse にHTTP Responseの状態を渡すためのstruct
type HTTPResponse struct {
	Status int
}

func FlushWithHTTPResponse(ctx context.Context, response *HTTPResponse) {
	l, ok := Value(ctx)
	if !ok {
		log.Print("failed FlushWithHTTPResponse. Logging.context not include LogContainer\n")
		return
	}
	if l.Entry.HTTPRequest == nil {
		log.Print("failed FlushWithHTTPResponse. HttpRequest is nil\n")
		return
	}
	l.Entry.HTTPRequest.Status = response.Status

	j, err := json.Marshal(l.Entry)
	if err != nil {
		log.Printf("failed LogContainer to json.err=%+v,LogContainer=%+v\n", err, l)
	}
	fmt.Println(string(j))
}

func SetTraceID(ctx context.Context, traceID string) {
	l, ok := Value(ctx)
	if !ok {
		log.Print("failed SetTraceID Logging.context not include LogContainer\n")
		return
	}

	l.Entry.TraceID = traceID
}

func SetSpanID(ctx context.Context, spanID string) {
	l, ok := Value(ctx)
	if !ok {
		log.Print("failed SetSpanID Logging.context not include LogContainer\n")
		return
	}

	l.Entry.SpanID = spanID
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
	l.Entry.Messages = append(l.Entry.Messages, string(j))
}

func createLogContainer(now time.Time) *LogContainer {
	project, err := gcpmetadata.GetProjectID()
	if err != nil {
		panic(err)
	}
	return &LogContainer{
		Entry: LogEntry{
			LogName: fmt.Sprintf("projects/%s/logs/slog", project),
			Timestamp: Timestamp{
				Seconds: now.Unix(),
				Nanos:   now.Nanosecond(),
			},
		},
	}
}
