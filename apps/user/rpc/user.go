package main

import (
	"flag"
	"fmt"
	"sync"

	"imooc.com/easy-chat/pkg/configserver"
	"imooc.com/easy-chat/pkg/interceptor/rpcserver"

	"imooc.com/easy-chat/apps/user/rpc/internal/config"
	"imooc.com/easy-chat/apps/user/rpc/internal/server"
	"imooc.com/easy-chat/apps/user/rpc/internal/svc"
	"imooc.com/easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

var grpcServer *grpc.Server
var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config

	err := configserver.NewConfigServer(*configFile, configserver.NewNacos(&configserver.NacosConfig{
		Addr:      "localhost",
		Namespace: "32d1730b-5d98-4cde-b0eb-52fcda17e2f8", // user 命名空间的实际 ID
		Group:     "DEFAULT_GROUP",
		DataId:    "user-rpc.yaml",
		Username:  "nacos",
		Password:  "wmlhy15723351559",
		LogLevel:  "warn",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		grpcServer.GracefulStop()

		wg.Add(1)
		go func(nc config.Config) {
			defer wg.Done()

			Run(c)
		}(c)
		return nil
	})

	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()

		Run(c)

	}(c)

	wg.Wait()
}

func Run(c config.Config) {
	ctx := svc.NewServiceContext(c)

	if err := ctx.SetRootToken(); err != nil {
		panic(err)
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(srv *grpc.Server) {
		fmt.Println("etcd = ", c.RpcServerConf)
		grpcServer = srv

		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
