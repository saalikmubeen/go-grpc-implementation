package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/saalikmubeen/go-grpc-implementation/authToken"
	generated_db "github.com/saalikmubeen/go-grpc-implementation/db/sqlc"
	"github.com/saalikmubeen/go-grpc-implementation/utils"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config    utils.Config
	store     generated_db.Store
	router    *gin.Engine
	authToken authToken.Maker
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config utils.Config, store generated_db.Store) (*Server, error) {

	tokenMaker, err := authToken.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)

	}

	server := &Server{store: store, authToken: tokenMaker, config: config}
	router := gin.Default()

	// Register a new CUSTOM validation function for currency with Gin.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.authToken))

	authRoutes.GET("/accounts", server.listUserAccounts)
	authRoutes.GET("/accounts/all", server.listAllAccounts)
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router

	return server, nil
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H { // gin.H -> map[string]interface{}
	return gin.H{"error": err.Error()}
}
