/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/zeromicro/go-zero/core/service"
	"imooc.com/easy-chat/apps/task/mq/internal/config"
	"imooc.com/easy-chat/apps/task/mq/internal/handler"
	"imooc.com/easy-chat/apps/task/mq/internal/svc"
	"imooc.com/easy-chat/pkg/configserver"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")
var wg sync.WaitGroup
var serviceGroup *service.ServiceGroup

func main() {
	flag.Parse()

	var c config.Config
	err := configserver.NewConfigServer(*configFile, configserver.NewNacos(&configserver.NacosConfig{
		Addr:      "127.0.0.1",
		Namespace: "task",
		Group:     "DEFAULT_GROUP",
		DataId:    "task-mq.yaml",
		Username:  "nacos",
		Password:  "nacos",
		LogLevel:  "warn",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		// 配置更新时，停止旧服务并重启
		if serviceGroup != nil {
			serviceGroup.Stop()
		}

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
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting mqueue at ...")
	serviceGroup.Start()
}
