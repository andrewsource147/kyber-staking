package api

import (
	"context"
	"fmt"
	_ "github.com/kyber/staking/api/docs"
	"github.com/kyber/staking/api/handler"
	"github.com/kyber/staking/contestant"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	"github.com/swaggo/echo-swagger"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	HTTP_Port = "8080"
)

type Server struct {
	Router  *echo.Echo
	Handler *handler.Handler
	mu      *sync.RWMutex
}

// @title Kyber Staking API
// @version 1.0
// @description API is used to interact with Kyber Staking contract
// @license.name Apache 2.0
// @host localhost:8080
func NewServer(kyberSc *contestant.KyberStakingContract) *Server {
	e := echo.New()

	e.Logger.SetLevel(echoLog.DEBUG)

	// Init middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Add validator
	e.Validator = NewValidator()

	h := handler.NewHandler(kyberSc)

	return &Server{
		Router:  e,
		Handler: h,
		mu:      &sync.RWMutex{},
	}
}

func (self *Server) RegisterRoutes() {
	self.mu.Lock()
	e := self.Router
	h := self.Handler
	defer self.mu.Unlock()

	e.GET("/", h.GetIndex)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/newContract", h.CreateNewContract)
	e.POST("/stake", h.Stake)
	e.POST("/withdraw", h.Withdraw)
	e.POST("/delegate", h.Delegate)
	e.POST("/vote", h.Vote)

	e.GET("/stake", h.GetStake)
	e.GET("/delegateStake", h.GetDelegateStake)
	e.GET("/representative", h.GetRepresentative)
	e.GET("/reward", h.GetReward)
	e.GET("/poolReward", h.GetPoolReward)
}

func (self *Server) Run() {
	self.RegisterRoutes()
	go func() {
		if err := self.Router.Start(fmt.Sprintf(":%s", HTTP_Port)); err != nil {
			log.Println("Shutting do¡wn the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := self.Router.Shutdown(ctx); err != nil {
		self.Router.Logger.Fatal(err)
	}
}
