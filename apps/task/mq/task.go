/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"imooc.com/easy-chat/apps/task/mq/internal/config"
	"imooc.com/easy-chat/apps/task/mq/internal/handler"
	"imooc.com/easy-chat/apps/task/mq/internal/svc"
	"sync"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")
var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	Run(c)
	//err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
	//	ETCDEndpoints:  "192.168.117.24:3379",
	//	ProjectKey:     "2f5bb7747efda0546636fb385a3fa593",
	//	Namespace:      "task",
	//	Configs:        "task-mq.yaml",
	//	ConfigFilePath: "./etc/conf",
	//	LogLevel:       "DEBUG",
	//})).MustLoad(&c, func(bytes []byte) error {
	//	var c config.Config
	//	configserver.LoadFromJsonBytes(bytes, &c)
	//
	//	wg.Add(1)
	//	go func(c config.Config) {
	//		defer wg.Done()
	//
	//		Run(c)
	//	}(c)
	//	return nil
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//wg.Add(1)
	//go func(c config.Config) {
	//	defer wg.Done()
	//
	//	Run(c)
	//}(c)
	//
	//wg.Wait()
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
