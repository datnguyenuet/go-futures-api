package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	swaggerfiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	authUseCase "go-futures-api/internal/auth/usecase"
	"go-futures-api/internal/interceptors"
	"go-futures-api/internal/middleware"
	positionsHandler "go-futures-api/internal/positions/delivery/http/v1"
	"go-futures-api/pkg/grpc_client"
	userService "go-futures-api/proto/auth"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	gin    *gin.Engine
	logger *zap.SugaredLogger
	tracer opentracing.Tracer
}

func NewServer(logger *zap.SugaredLogger, tracer opentracing.Tracer) *Server {
	return &Server{
		gin:    gin.Default(),
		logger: logger,
		tracer: tracer,
	}
}

func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	im := interceptors.NewInterceptorManager(s.logger, s.tracer)

	validate := validator.New()

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", "5555")
		if err := http.ListenAndServe("localhost:5555", http.DefaultServeMux); err != nil {
			s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
		}
	}()

	authServiceConn, err := grpc_client.NewGRPCClientServiceConn(ctx, im, ":5005")
	authServiceClient := userService.NewUserServiceClient(authServiceConn)
	authUC := authUseCase.NewAuthUseCase(s.logger, authServiceClient)

	mw := middleware.NewMiddlewareManager(s.logger, authUC)

	v1 := s.gin.Group("/api/v1")
	ordersGroup := v1.Group("/positions")
	positionsGroup := v1.Group("/positions")
	positionsGroup.Use(mw.JWTMiddleware())

	positionHandlers := positionsHandler.NewPositionsHandlers(positionsGroup, s.logger)
	positionHandlers.MapRoutes()

	//commentHandlers := commentsHandlers.NewCommentsHandlers(s.cfg, commentsGroup, s.logger, validate, commUC, mw)
	//commentHandlers.MapRoutes()

	s.MapRoutes()
	err := s.gin.Run(":5555")
	if err != nil {
		return err
	}

	fmt.Println(ctx, im, validate, ordersGroup, positionsGroup)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.logger.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.logger.Errorf("ctx.Done: %v", done)
	}
	return nil
}

func (s *Server) MapRoutes() {
	err := s.gin.SetTrustedProxies([]string{"172.16.0.131"})
	if err != nil {
		return
	}
	s.gin.TrustedPlatform = gin.PlatformGoogleAppEngine
	// Or set your own trusted request header for another trusted proxy service
	// Don't set it to any suspect request header, it's unsafe
	s.gin.TrustedPlatform = "X-CDN-IP"

	s.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	s.gin.GET("/health", func(cxt *gin.Context) {
		cxt.String(http.StatusOK, "Ok")
	})
}
