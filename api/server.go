package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/taufiqdp/go-simplebank/db/sqlc"
)

type Server struct {
	db     db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{db: store}
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
	return server
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"Error": err.Error()}
}
