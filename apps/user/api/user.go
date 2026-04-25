package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"imooc.com/easy-chat/apps/user/api/internal/config"
	"imooc.com/easy-chat/apps/user/api/internal/handler"
	"imooc.com/easy-chat/apps/user/api/internal/svc"
	"imooc.com/easy-chat/pkg/configserver"
	"imooc.com/easy-chat/pkg/resultx"
)

var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	//conf.MustLoad(*configFile, &c)

	err := configserver.NewConfigServer(*configFile, configserver.NewNacos(&configserver.NacosConfig{
		Addr:      "127.0.0.1",
		Namespace: "user",
		Group:     "DEFAULT_GROUP",
		DataId:    "user-api.yaml",
		Username:  "nacos",
		Password:  "nacos",
		LogLevel:  "warn",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		proc.WrapUp()

		wg.Add(1)
		go func(c config.Config) {
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
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OkHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
