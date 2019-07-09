package v2_test

import (
	"context"
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

func TestInfo(t *testing.T) {
	ctx := slog.WithValue(context.Background())
	slog.Info(ctx, slog.KV{"hoge", "fuga"})
	lc, ok := slog.Value(ctx)
	if !ok {
		t.Error("want ok but ng...")
	}

	if e, g := 1, len(lc.Messages); e != g {
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
	// Output: {"entry":{"timestamp":{"seconds":1546398245,"nanos":6},"message":"[\"{\\\"key\\\":\\\"hoge\\\",\\\"value\\\":\\\"fuga\\\"}\"]","severity":"","thread":1546398245000000006,"httpRequest":{"requestUrl":""}},"messages":["{\"key\":\"hoge\",\"value\":\"fuga\"}"]}
}
