package slog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Entry is Stackdriver Logging Entry
type Entry struct {
	LogName     string            `json:"logName"`
	Resource    MonitoredResource `json:"resource"`
	JSONPayload interface{}       `json:"jsonPayload"`
}

// MonitoredResource is Log Resource
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/MonitoredResource
type MonitoredResource struct {
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

// Log is Log Object
type Log struct {
	Entry    Entry `json:"entry"`
	Messages []string
}

var logMap map[context.Context]*Log

var m sync.RWMutex

func init() {
	logMap = make(map[context.Context]*Log)
	m = sync.RWMutex{}
}

func setLogMap(ctx context.Context, log *Log) {
	m.Lock()
	defer m.Unlock()
	logMap[ctx] = log
}

func getLogMap(ctx context.Context) (*Log, bool) {
	m.RLock()
	defer m.RUnlock()
	e, ok := logMap[ctx]
	return e, ok
}

// Info is output info level Log
func Info(ctx context.Context, message string) {
	e, ok := getLogMap(ctx)
	if !ok {
		e = &Log{
			Entry: Entry{
				LogName: "projects/metal-tile-dev1/logs/slog",
				Resource: MonitoredResource{
					Type: "slog",
				},
				//Severity:  "INFO",
				//Timestamp: time.Now(),
			},
		}
		go log(ctx)
	}
	e.Messages = append(e.Messages, message)
	setLogMap(ctx, e)
}

func log(ctx context.Context) {
	fmt.Println("log start")
	select {
	case <-ctx.Done():
		fmt.Println("log ctx.Done()")
		e, ok := logMap[ctx]
		if ok {
			fmt.Println("log ctx.Done() logMap ok;")
			encoder := json.NewEncoder(os.Stdout)
			e.Entry.JSONPayload = e.Messages
			if err := encoder.Encode(e.Entry); err != nil {
				_, err := os.Stderr.WriteString(err.Error())
				if err != nil {
					panic(err)
				}
			}
			fmt.Println("log ctx.Done() logMap ok; done.")
		}
	}
	fmt.Println("log end")
}
