package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/poc-git/semrelease"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger, _ := configLog(zap.NewAtomicLevelAt(zapcore.InfoLevel)).Build()
	defer logger.Sync()
	repository := semrelease.NewRepository()

	service := semrelease.NewService(repository)
	handler := createServerHandler(service, logger)
	gracehttp.Serve(&http.Server{
		Addr:         fmt.Sprintf(":%d", 9090),
		Handler:      handler,
		ReadTimeout:  time.Duration(10 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
	})

}
func configLog(level zap.AtomicLevel) zap.Config {
	return zap.Config{
		Level:         level,
		Development:   false,
		DisableCaller: true,
		Sampling:      nil,
		Encoding:      "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: millisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func millisDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt(int(float64(d) / float64(time.Millisecond)))
}
