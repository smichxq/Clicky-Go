package mtl

import (
	"io"
	"os"
	"time"

	kitexzap "github.com/kitex-contrib/obs-opentelemetry/logging/zap"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog(ioWriter io.Writer) {
	var opts []kitexzap.Option
	var output zapcore.WriteSyncer
	if os.Getenv("GO_ENV") != "online" {
		opts = append(opts, kitexzap.WithCoreEnc(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())))
		output = zapcore.AddSync(ioWriter)
	} else {
		opts = append(opts, kitexzap.WithCoreEnc(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())))
		// async log
		output = &zapcore.BufferedWriteSyncer{
			WS:            zapcore.AddSync(ioWriter),
			FlushInterval: time.Minute,
		}
	}
	server.RegisterShutdownHook(func() {
		output.Sync() //nolint:errcheck
	})
	log := kitexzap.NewLogger(opts...)
	klog.SetLogger(log)
	klog.SetLevel(klog.LevelTrace)
	klog.SetOutput(output)
}
