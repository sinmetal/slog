package slog

import (
	"context"
	"fmt"
	"testing"
)

func Example() {
	//if e, g := 1, len(log.Messages); e != g {
	//	t.Fatalf("log.messages.len expected %d; got %d", e, g)
	//}
	//
	//if e, g := "Hello slog World", log.Messages[0]; e != g {
	//	t.Fatalf("log.messages[0] expected %s; got %s", e, g)
	//}
}

func TestLog_Info(t *testing.T) {
	//loc, err := time.LoadLocation("Asia/Tokyo")
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	{
		handleLog("First")
	}
	{
		handleLog("Second")
	}

	// Output: {"insertId":"2018-06-12 09:08:58.07136086 +0900 JST m=+0.000586925","severity":"WARNING","labels":{"hoge":"fuga"},"logName":"projects/metal-tile-dev1/logs/slog","receiveTimestamp":"2018-06-12T09:08:58.071473299+09:00","resource":{"type":"slog","labels":{"hoge":"fuga"}},"jsonPayload":["Hello First 1","Hello First 2"],"timestamp":"2018-06-12T09:08:58.07147335+09:00"}
}

// TestLog_Empty is ログが空っぽの場合は出力しないことを確認
func TestLog_Empty(t *testing.T) {
	ctx := context.Background()
	ctx = WithLog(ctx)
	defer Flush(ctx)
}

// TestLogLevel is LogLevelがWarningになることを確認
func TestLogLevel(t *testing.T) {
	ctx := context.Background()
	ctx = WithLog(ctx)
	defer Flush(ctx)

	SetLogName(ctx, "Test Log Level")
	Warning(ctx, "Hoge", "Fuga")
	Info(ctx, "Hoge", "Fuga")

	l, ok := ctx.Value(contextLogKey{}).(*StackdriverLogEntry)
	if !ok {
		t.Errorf("failed get Log")
	}
	if e, g := "WARNING", l.Severity; e != g {
		t.Errorf("expected Severity is %s; got %s", e, g)
	}
}

func handleLog(message string) {
	ctx := context.Background()
	ctx = WithLog(ctx)
	defer Flush(ctx)
	Info(ctx, "handleLog", KV{"message", fmt.Sprintf("Hello %s 1", message)})
	Info(ctx, "handleLog", fmt.Sprintf("Hello %s 2", message))
	Info(ctx, "handleLog", KV{"message", 3})
}

func TestLog_InfoWithCancel(t *testing.T) {
	//loc, err := time.LoadLocation("Asia/Tokyo")
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	{
		handleLogWithCancel("First")
	}
	{
		handleLogWithCancel("Second")
	}

	// Output: {"insertId":"2018-06-12 09:08:58.07136086 +0900 JST m=+0.000586925","severity":"WARNING","labels":{"hoge":"fuga"},"logName":"projects/metal-tile-dev1/logs/slog","receiveTimestamp":"2018-06-12T09:08:58.071473299+09:00","resource":{"type":"slog","labels":{"hoge":"fuga"}},"jsonPayload":["Hello First 1","Hello First 2"],"timestamp":"2018-06-12T09:08:58.07147335+09:00"}
}

func handleLogWithCancel(message string) {
	ctx := context.Background()
	ctx = WithLog(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer Flush(ctx)
	Info(ctx, "handleLogWithCancel", fmt.Sprintf("Hello WithCancel %s 1", message))
	cancel()
	Info(ctx, "handleLogWithCancel", fmt.Sprintf("Hello WithCancel %s 2", message))
}

//func ExampleLog_Infof() {
//	loc, err := time.LoadLocation("Asia/Tokyo")
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	log := Start(time.Date(2017, time.April, 1, 13, 15, 30, 45, loc))
//	log.Infof("Hello World %d", 1)
//	log.Infof("Hello World %d", 2)
//	log.Flush()
//	// Output: {"timestamp":{"seconds":1491020130,"nanos":45},"message":"[\"Hello World 1\",\"Hello World 2\"]","severity":"INFO","thread":1491020130000000045}
//}
//
//func ExampleLog_Info() {
//	loc, err := time.LoadLocation("Asia/Tokyo")
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	log := Start(time.Date(2017, time.April, 1, 13, 15, 30, 45, loc))
//	log.Info("Hello World 1")
//	log.Info("Hello World 2")
//	log.Flush()
//	// Output: {"timestamp":{"seconds":1491020130,"nanos":45},"message":"[\"Hello World 1\",\"Hello World 2\"]","severity":"INFO","thread":1491020130000000045}
//}
//
//func ExampleLog_Errorf() {
//	loc, err := time.LoadLocation("Asia/Tokyo")
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	log := Start(time.Date(2017, time.April, 1, 13, 15, 30, 45, loc))
//	log.Errorf("Hello World %d", 1)
//	log.Errorf("Hello World %d", 2)
//	log.Flush()
//	// Output: {"timestamp":{"seconds":1491020130,"nanos":45},"message":"[\"Hello World 1\",\"Hello World 2\"]","severity":"ERROR","thread":1491020130000000045}
//}
//
//func ExampleLog_Error() {
//	loc, err := time.LoadLocation("Asia/Tokyo")
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	log := Start(time.Date(2017, time.April, 1, 13, 15, 30, 45, loc))
//	log.Error("Hello World 1")
//	log.Error("Hello World 2")
//	log.Flush()
//	// Output: {"timestamp":{"seconds":1491020130,"nanos":45},"message":"[\"Hello World 1\",\"Hello World 2\"]","severity":"ERROR","thread":1491020130000000045}
//}
//
//// ExampleLog_Error2 is Error Levelが優先されることを確認
//func ExampleLog_Error_level() {
//	loc, err := time.LoadLocation("Asia/Tokyo")
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	log := Start(time.Date(2017, time.April, 1, 13, 15, 30, 45, loc))
//	log.Info("Hello Info")
//	log.Error("Hello Error")
//	log.Info("Hello Info")
//	log.Flush()
//	// Output: {"timestamp":{"seconds":1491020130,"nanos":45},"message":"[\"Hello Info\",\"Hello Error\",\"Hello Info\"]","severity":"ERROR","thread":1491020130000000045}
//}

func Test_maxSeverity(t *testing.T) {
	if e, g := "WARNING", maxSeverity("INFO", "WARNING"); e != g {
		t.Errorf("expected Severity is %s; got %s", e, g)
	}
	if e, g := "WARNING", maxSeverity("WARNING", "INFO"); e != g {
		t.Errorf("expected Severity is %s; got %s", e, g)
	}
}
