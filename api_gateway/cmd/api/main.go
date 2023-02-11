package main

import (
	"github.com/opentracing/opentracing-go"
	_ "go-futures-api/docs"
	"go-futures-api/internal/server"
	"go-futures-api/pkg/jaeger"
	"go.uber.org/zap"
)

// @title API Gateway
// @version 1.0
// @description API Gateway
// @contact.name futures
// @contact.url https://google.com/
// @contact.email futures@mail.com
// @BasePath /api/v1
func main() {
	logger := zap.S()
	tracer, closer, err := jaeger.InitJaeger()
	if err != nil {
		logger.Fatal("cannot create tracer", err)
	}

	logger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	logger.Info("Opentracing connected")

	s := server.NewServer(logger, tracer)
	logger.Fatal(s.Run())
}
