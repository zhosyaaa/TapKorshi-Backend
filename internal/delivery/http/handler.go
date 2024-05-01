package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zhosyaaa/RoommateTap/internal/config"
	"github.com/zhosyaaa/RoommateTap/internal/delivery/http/v1"
	"github.com/zhosyaaa/RoommateTap/internal/service"
	"github.com/zhosyaaa/RoommateTap/pkg/auth"
	"github.com/zhosyaaa/RoommateTap/pkg/limiter"
	"golang.org/x/oauth2"
	"net/http"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
	cfg          oauth2.Config
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager, cfg oauth2.Config) *Handler {
	return &Handler{services: services, tokenManager: tokenManager, cfg: cfg}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL),
		corsMiddleware,
	)

	//docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	//if cfg.Environment != config.EnvLocal {
	//	docs.SwaggerInfo.Host = cfg.HTTP.Host
	//}

	//if cfg.Environment != config.Prod {
	//	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//}

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager, h.cfg)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
