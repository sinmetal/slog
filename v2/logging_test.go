package v2_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/sinmetal/silverdog/dogtime"
	slog "github.com/sinmetal/slog/v2"
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

	r := http.Request{
		Method: "GET",
		URL: &url.URL{
			Host:    "sinmetal.jp",
			Path:    "/api/v1/hoge",
			RawPath: "/api/v1/hoge",
		},
		Proto: "HTTP 1.1",
	}

	ctx := slog.WithValueForHTTP(context.Background(), r)
	_, ok := slog.Value(ctx)
	if !ok {
		return
	}
	slog.Info(ctx, slog.KV{"hoge", "fuga"})

	slog.FlushWithHTTPResponse(ctx, 200)
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
