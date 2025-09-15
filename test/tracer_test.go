package main

import (
	"context"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/zrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
	"imooc.com/easy-chat/apps/user/rpc/userclient"
	"testing"
)

var tracerTestAttributeKey1 = attribute.Key("test.tracer.desc")

//func StartAgent(name string) {
//	trace.StartAgent(trace.Config{
//		Name:     name,
//		Endpoint: "http://192.168.117.80:14268/api/traces",
//		Batcher:  "jaeger",
//		Disabled: false,
//		Sampler:  1,
//	})
//}

func Test_Tracer(t *testing.T) {
	// 对tracer初始化

	// 执行第一个反复

	t.Run("启动", func(t *testing.T) {
		StartAgent("test.go-zero.tracer")
		ctx, span := otel.Tracer(trace.TraceName).Start(context.Background(), "a")

		span.SetStatus(codes.Ok, "")
		t.Log("isRecording", span.IsRecording())
		exec(t, ctx, "exec")
		span.End()
		for {
		}
	})

	t.Run("加载", func(t *testing.T) {
		StartAgent("test.go-zero.tracer--2")
		tr := otel.Tracer(trace.TraceName)

		// Get the trace ID and span ID from the previous span

		previousTraceID, err := oteltrace.TraceIDFromHex("1854784d177d6801e0ba2c96f2620eea")
		if err != nil {
			t.Error("oteltrace.TraceIDFromHex err ", err)
		}
		previousSpanID, err := oteltrace.SpanIDFromHex("4c6310df3367138a")
		if err != nil {
			t.Error("oteltrace.TraceIDFromHex err ", err)
		}

		// Create a new span context based on the previous trace ID and span ID
		spanContext := oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID:    previousTraceID,
			SpanID:     previousSpanID,
			TraceFlags: oteltrace.FlagsSampled,
		})

		// Start a new span with the previous span context
		ctx, span := tr.Start(oteltrace.ContextWithRemoteSpanContext(context.Background(), spanContext), "example-span")
		span.SetStatus(codes.Ok, "")
		t.Log("isRecording", span.IsRecording())
		exec(t, ctx, "加载")
		span.End()

		for {
		}
	})

	// 结合go-zero中的rpc请求
	t.Run("发送", func(t *testing.T) {

		StartAgent("test.go-zero.tracer--3")
		tr := otel.Tracer(trace.TraceName)

		// Get the trace ID and span ID from the previous span

		previousTraceID, err := oteltrace.TraceIDFromHex("1854784d177d6801e0ba2c96f2620eea")
		if err != nil {
			t.Error("oteltrace.TraceIDFromHex err ", err)
		}
		previousSpanID, err := oteltrace.SpanIDFromHex("4c6310df3367138a")
		if err != nil {
			t.Error("oteltrace.TraceIDFromHex err ", err)
		}
		spanContext := oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID:    previousTraceID,
			SpanID:     previousSpanID,
			TraceFlags: oteltrace.FlagsSampled,
		})
		ctx, span := tr.Start(oteltrace.ContextWithRemoteSpanContext(context.Background(), spanContext), "example-span")
		span.SetStatus(codes.Ok, "")
		userRpc := userclient.NewUser(zrpc.MustNewClient(zrpc.RpcClientConf{
			Etcd: discov.EtcdConf{
				Hosts: []string{"192.168.117.80:3379"},
				Key:   "user.rpc",
			},
			Timeout: 2000,
			Middlewares: zrpc.ClientMiddlewaresConf{
				Trace:      true,
				Duration:   true,
				Prometheus: true,
				Breaker:    true,
				Timeout:    true,
			},
		}))
		resp, err := userRpc.FindUser(ctx, &userclient.FindUserReq{
			Name:  "",
			Phone: "17388880002",
			Ids:   nil,
		})
		t.Log(resp, err)
		span.End()

		for {
		}
	})

}

func exec(t *testing.T, ctx context.Context, desc string) {
	_, span := startSpan(ctx, "exec", desc)
	defer span.End()

	t.Log("trace id ", span.SpanContext().TraceID().String())
	t.Log("span id ", span.SpanContext().SpanID().String())
	t.Log("执行任务 .... ", desc)
}

func startSpan(ctx context.Context, spanName, desc string) (context.Context, oteltrace.Span) {
	// 该方法会去获取到tracer，如果存在
	tracer := trace.TracerFromContext(ctx)
	start, span := tracer.Start(ctx, spanName)
	span.SetAttributes(tracerTestAttributeKey.String(desc))
	span.SetStatus(codes.Ok, "")
	return start, span
}
