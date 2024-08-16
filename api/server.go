package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/taufiqdp/go-simplebank/db/sqlc"
	"github.com/taufiqdp/go-simplebank/token"
	"github.com/taufiqdp/go-simplebank/utils"
)

type Server struct {
	config     utils.Config
	db         db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		db:         store,
		tokenMaker: tokenMaker,
	}

	gin.SetMode(gin.DebugMode)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	engine := gin.Default()
	engine.SetTrustedProxies([]string{"127.0.0.1"})

	engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	engine.POST("/users", server.CreateUser)
	engine.GET("/users/:username", server.GetUser)

	engine.POST("/accounts", server.CreateAccount)
	engine.GET("/accounts/:id", server.GetAccount)
	engine.GET("/accounts", server.ListAccount)

	engine.POST("/transfers", server.CreateTransfer)

	server.router = engine
	return server, nil
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"Error": err.Error()}
}
