package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/modules/middlewares/middlewareHandler"
	"github.com/guatom999/Go-MicroService/modules/middlewares/middlewareRepositories"
	"github.com/guatom999/Go-MicroService/modules/middlewares/middlewareUseCases"
	"github.com/guatom999/Go-MicroService/pkg/jwtauth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	server struct {
		app        *echo.Echo
		db         *mongo.Client
		cfg        *config.Config
		middleware middlewareHandler.IMiddlewareHandlerService
	}
)

func (s *server) gracefulShutdown(pctx context.Context, close <-chan os.Signal) {

	log.Printf("Start service: %s", s.cfg.App.Name)

	<-close

	log.Printf("Shutting down service: %s", s.cfg.App.Name)

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	if err := s.app.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown:%v", err)
	}
}

func (s *server) httpListening() {

	if err := s.app.Start(s.cfg.App.Url); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to shutdown:%v", err)
	}

}

func NewMiddleware(cfg *config.Config) middlewareHandler.IMiddlewareHandlerService {
	repository := middlewareRepositories.NewMiddlewareRepository()
	usecase := middlewareUseCases.NewMiddlewareUseCase(repository)
	handler := middlewareHandler.NewMiddlewareUseCaseService(cfg, usecase)

	return handler
}

func Start(pctx context.Context, cfg *config.Config, db *mongo.Client) {
	s := &server{
		app:        echo.New(),
		db:         db,
		cfg:        cfg,
		middleware: NewMiddleware(cfg),
	}

	jwtauth.SetApiKey(cfg.Jwt.ApiSecretKey)

	// Request Timeout
	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      time.Second * 10,
	}))

	//Cors
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	// Body Limit
	s.app.Use(middleware.BodyLimit("10M"))

	switch s.cfg.App.Name {
	case "auth":
		s.authService()
	case "inventory":
		s.inventoryService()
	case "item":
		s.itemService()
	case "payment":
		s.paymentService()
	case "player":
		s.playerService()
	}

	//close
	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM)

	s.app.Use(middleware.Logger())

	go s.gracefulShutdown(pctx, close)

	// go func() {
	// 	_ = <-c
	// 	log.Println("server is shutting down....")
	// 	_ = s.app.Shutdown()
	// }()

	s.httpListening()

}
