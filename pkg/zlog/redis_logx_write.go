package zlog

import (
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

const (
	levelAlert  = "alert"
	levelInfo   = "info"
	levelError  = "error"
	levelSevere = "severe"
	levelFatal  = "fatal"
	levelSlow   = "slow"
	levelStat   = "stat"
	levelDebug  = "debug"
)

const (
	callerKey    = "caller"
	contentKey   = "content"
	durationKey  = "duration"
	levelKey     = "level"
	spanKey      = "span"
	timestampKey = "@timestamp"
	traceKey     = "trace"
	truncatedKey = "truncated"
)

type redisLogxWriter struct {
	redisKey string
	redis    *redis.Redis
}

func NewRedisLogxWriter(redisKey string, redisCfg redis.RedisConf) logx.Writer {
	return &redisLogxWriter{
		redisKey: redisKey,
		redis:    redis.MustNewRedis(redisCfg),
	}
}

func (n *redisLogxWriter) Alert(v any) {
	n.outPut(v, levelAlert)
}

func (n *redisLogxWriter) Close() error {
	return nil
}

func (n *redisLogxWriter) Debug(v any, fields ...logx.LogField) {
	n.outPut(v, levelDebug, fields...)
}

func (n *redisLogxWriter) Error(v any, fields ...logx.LogField) {
	n.outPut(v, levelError, fields...)
}

func (n *redisLogxWriter) Info(v any, fields ...logx.LogField) {
	n.outPut(v, levelInfo, fields...)
}

func (n *redisLogxWriter) Severe(v any) {
	n.outPut(v, levelSevere)
}

func (n *redisLogxWriter) Slow(v any, fields ...logx.LogField) {
	n.outPut(v, levelSlow, fields...)
}

func (n *redisLogxWriter) Stack(v any) {
	n.outPut(v, levelError)
}

func (n *redisLogxWriter) Stat(v any, fields ...logx.LogField) {
	n.outPut(v, levelStat, fields...)
}

func (n *redisLogxWriter) outPut(v any, level string, fields ...logx.LogField) {
	// 根据日志处理数据信息

	// 日志内容
	content := make(map[string]any)

	content[contentKey] = v
	for _, field := range fields {
		t := logx.Field(field.Key, field.Value)
		content[t.Key] = t.Value
	}

	// 增加时间
	content[timestampKey] = time.Now().UnixNano()
	// 等级
	content[levelKey] = level

	f, ok := getCallerFrame(6)
	if ok {
		content[callerKey] = fmt.Sprintf("%v:%v", f.File, f.Line)
	}

	// 格式化

	b, err := json.Marshal(content)
	if err != nil {
		return
	}

	data := string(b)

	// 输出
	go n.redis.Rpush(n.redisKey, data)
}
