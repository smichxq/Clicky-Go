package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"clicky.website/clicky/security/biz/dal"
	"clicky.website/clicky/security/conf"
	"clicky.website/clicky/security/kitex_gen/security/security"
	"clicky.wesite/clicky/common/mtl"
	"clicky.wesite/clicky/common/serversuite"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	ServiceName      = conf.GetConf().Kitex.Service
	RegisterAddr     = conf.GetConf().Registry.RegistryAddress[0]
	ConsulHealthAddr = conf.GetConf().Kitex.ConsulHealthAddr
	MetricsPort      = conf.GetConf().Kitex.MetricsPort
)

func main() {
	optl := mtl.InitTracing(ServiceName)
	defer optl.Shutdown(context.Background())

	mtl.InitMetric(ServiceName, MetricsPort, RegisterAddr)

	dal.Init()

	opts := kitexInit()

	// health check
	go StartHealthCheckServer(":" + strings.Split(ConsulHealthAddr, ":")[1])

	svr := security.NewServer(new(SecurityImpl), opts...)

	err := svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	suite := serversuite.CommonServerSuite{
		CurrentServiceName: ServiceName,
		RegistryAddr:       RegisterAddr,
		ConsulHealthAddr:   ConsulHealthAddr,
	}

	opts = append(opts, suite.Options()...)
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))
	// thrift meta handler
	opts = append(opts, server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	// klog
	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Kitex.LogFileName,
			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		}),
		FlushInterval: time.Minute,
	}
	klog.SetOutput(asyncWriter)
	server.RegisterShutdownHook(func() {
		asyncWriter.Sync()
	})
	return
}

// 健康监测接口
func StartHealthCheckServer(addr string) {
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("health check server error: %v", err)
		}
	}()
}
