package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/taufiqdp/go-simplebank/db/sqlc"
)

type Server struct {
	db     *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{db: store}
	gin.SetMode(gin.DebugMode)

	engine := gin.New()
	engine.SetTrustedProxies([]string{"127.0.0.1"})

	engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	engine.POST("/accounts", server.CreateAccount)
	engine.GET("/accounts/:id", server.GetAccount)
	engine.GET("/accounts", server.ListAccount)

	server.router = engine
	return server
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"Error": err.Error()}
}
