package handlers

import (
	"net/http"

	"github.com/avalonprod/gasstrem/src/internal/config"
	"github.com/avalonprod/gasstrem/src/internal/services"
	"github.com/avalonprod/gasstrem/src/packages/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	services     *services.Services
	tokenManager auth.TokenManager
}

func NewHandlers(services *services.Services, tokenManager auth.TokenManager) *Handlers {
	return &Handlers{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handlers) Init(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(
		gin.Recovery(),
		gin.Logger(),
		// h.CorsMiddleware,
		cors.Default(),
	)

	api := r.Group("/api")
	{
		authenticated := api.Group("/", h.userIndentity)
		{
			invoces := authenticated.Group("/invoces")
			{
				invoces.POST("/create", h.CreateInvoice)
				invoces.GET("/get-all", h.GetAllInvoceByUserId)
			}
		}

		users := api.Group("/users")
		{
			users.POST("/refresh", h.userRefreshToken)
			users.POST("/sign-up", h.UsersSignUp)
			users.POST("/sign-in", h.UsersSignIn)
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}
