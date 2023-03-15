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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	l := zap.S()
	tracer, closer, err := jaeger.InitJaeger()
	if err != nil {
		l.Fatal("cannot create tracer", err)
	}

	l.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	l.Info("Opentracing connected")

	s := server.NewServer(l, tracer)
	l.Fatal(s.Run())
}
