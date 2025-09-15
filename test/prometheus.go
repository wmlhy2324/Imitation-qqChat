package main

import (
	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/prometheus"
	"time"
)

func main() {
	// 初始化prometheus,并对外提供监听接口
	pCfg := prometheus.Config{
		Host: "0.0.0.0",
		Port: 1234,
		Path: "/metrics",
	}
	prometheus.StartAgent(pCfg)

	gaugeVec := metric.NewGaugeVec(&metric.GaugeVecOpts{
		Namespace: "core",
		Subsystem: "tests",
		Name:      "go_zero_test",
		Help:      "test go-zero prometheus and metric",
		Labels:    []string{"path"},
	})

	var i int
	for {
		i++
		if i%2 == 0 {
			gaugeVec.Inc("/user")
		}

		time.Sleep(time.Second)
	}
}

//func main() {
//
//	temp := prometheus.NewGaugeVec(prometheus.GaugeOpts{
//		Namespace: "go_zero",
//		Name:      "tests_temp_gauge",
//		Help:      "the is test gauge",
//	}, []string{"path"})
//	prometheus.MustRegister(temp)
//
//	var i int
//	go func() {
//		for {
//			i++
//			if i%2 == 0 {
//				temp.WithLabelValues("/user").Inc()
//			}
//
//			time.Sleep(time.Second)
//		}
//	}()
//
//	http.Handle("/metrics", promhttp.Handler())
//	fmt.Println("启动服务")
//	fmt.Println(http.ListenAndServe(":1234", nil))
//}
