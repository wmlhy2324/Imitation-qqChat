package zlog

import (
	"context"
	"gitee.com/dn-jinmin/tlog"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Tlog struct{}

func NewTlog(redisCfg redis.RedisConf) *Tlog {
	io := NewRedisIoWriter("www.imooc.com", redisCfg)
	logx.SetWriter(logx.NewWriter(io))
	return &Tlog{}
}

func (z *Tlog) Write(ctx context.Context, content *tlog.Content, tfields ...tlog.Field) {

	fields := z.buildFields(content, tfields...)

	switch content.Level {
	case tlog.INFO:
		logx.Infow(content.Msg, fields...)
	case tlog.DEBUG:
		logx.Debugw(content.Msg, fields...)
	case tlog.ERROR, tlog.FATAL:
		logx.Errorw(content.Msg, fields...)
	case tlog.SEVERE:
		logx.Severe(content.Msg)
	case tlog.ALERT:
		logx.Alert(content.Msg)
	case tlog.STAT:
		logx.Stat(content.Msg)
	case tlog.SLOW:
		logx.Sloww(content.Msg, fields...)
	}
}

func (z *Tlog) buildFields(content *tlog.Content, tfields ...tlog.Field) []logx.LogField {
	fields := make([]logx.LogField, 0, 4+len(tfields))
	fields = append(fields, logx.LogField{
		Key:   tlog.TlogLabel,
		Value: content.Label,
	})
	fields = append(fields, logx.LogField{
		Key:   tlog.TlogTraceId,
		Value: content.TraceId,
	})
	fields = append(fields, logx.LogField{
		Key:   tlog.TlogRelatedId,
		Value: content.RelatedId,
	})

	fields = append(fields, logx.LogField{
		Key:   tlog.TlogPath,
		Value: content.Path,
	})

	for _, tfield := range tfields {
		fields = append(fields, logx.LogField{
			Key:   tfield.Key,
			Value: tfield.Val,
		})
	}
	return fields
}
