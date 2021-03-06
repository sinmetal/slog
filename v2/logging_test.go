package slog_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/sinmetal/silverdog/dogtime"
	"github.com/sinmetal/slog/v2"
)

func TestWithValue(t *testing.T) {
	ctx := slog.WithValue(context.Background())
	_, ok := slog.Value(ctx)
	if !ok {
		t.Error("want ok but ng...")
	}
}

func ExampleForHTTP() {
	sn := dogtime.StockNower{}
	dogtime.SetNower(&sn)

	now := time.Date(2019, 1, 2, 3, 4, 5, 6, time.UTC)
	sn.AddStockTime(now)

	r := &http.Request{
		Method:     "GET",
		RequestURI: "/api/v1/hoge",
		Proto:      "HTTP 1.1",
	}

	ctx := slog.WithValueForHTTP(context.Background(), r)
	_, ok := slog.Value(ctx)
	if !ok {
		return
	}
	slog.Info(ctx, slog.KV{"hoge", "fuga"})

	var response slog.HTTPResponse
	response.Status = http.StatusOK
	slog.FlushWithHTTPResponse(ctx, &response)
	// Output: {"timestamp":{"seconds":1546398245,"nanos":6},"messages":["{\"key\":\"hoge\",\"value\":\"fuga\"}"],"httpRequest":{"requestMethod":"GET","requestUrl":"/api/v1/hoge","status":200,"protocol":"HTTP 1.1"}}
}

func TestInfo(t *testing.T) {
	ctx := slog.WithValue(context.Background())
	slog.Info(ctx, slog.KV{"hoge", "fuga"})
	lc, ok := slog.Value(ctx)
	if !ok {
		t.Error("want ok but ng...")
	}

	if e, g := 1, len(lc.Entry.Messages); e != g {
		t.Errorf("Messages.length want %v but got %v", e, g)
	}
}

func ExampleInfo() {
	sn := dogtime.StockNower{}
	dogtime.SetNower(&sn)

	now := time.Date(2019, 1, 2, 3, 4, 5, 6, time.UTC)
	sn.AddStockTime(now)

	ctx := slog.WithValue(context.Background())
	slog.Info(ctx, slog.KV{"hoge", "fuga"})
	slog.Flush(ctx)
	// Output: {"timestamp":{"seconds":1546398245,"nanos":6},"messages":["{\"key\":\"hoge\",\"value\":\"fuga\"}"]}
}
