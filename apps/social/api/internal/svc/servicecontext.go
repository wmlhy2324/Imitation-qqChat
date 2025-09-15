package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"imooc.com/easy-chat/apps/im/rpc/imclient"
	"imooc.com/easy-chat/apps/social/api/internal/config"
	"imooc.com/easy-chat/apps/social/rpc/socialclient"
	"imooc.com/easy-chat/apps/user/rpc/userclient"
	"imooc.com/easy-chat/pkg/middleware"
)

var retryPolicy = `{
	"methodConfig" : [{
		"name": [{
			"service": "social.social"
		}],
		"waitForReady": true,
		"retryPolicy": {
			"maxAttempts": 5,
			"initialBackoff": "0.001s",
			"maxBackoff": "0.002s",
			"backoffMultiplier": 1.0,
			"retryableStatusCodes": ["UNKNOWN", "DEADLINE_EXCEEDED"]
		}
	}]
}`

type ServiceContext struct {
	Config                config.Config
	IdempotenceMiddleware rest.Middleware
	LimitMiddleware       rest.Middleware
	*redis.Redis
	socialclient.Social
	userclient.User
	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		Redis:                 redis.MustNewRedis(c.Redisx),
		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware().Handler,
		LimitMiddleware:       middleware.NewLimitMiddleware(c.Redisx).TokenLimitHandler(1, 100),
		Social:                socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),//zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)),
		//zrpc.WithUnaryClientInterceptor(interceptor.DefaultIdempotentClient),

		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),//zrpc.WithUnaryClientInterceptor(rpcclient.NewSheddingClient("user-rpc",
		//	load.WithBuckets(10),
		//	load.WithCpuThreshold(1),
		//	load.WithWindow(time.Millisecond*100000),
		//)),

		Im: imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
	}
}
