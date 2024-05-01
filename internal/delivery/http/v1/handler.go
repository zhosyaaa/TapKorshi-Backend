package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhosyaaa/RoommateTap/internal/service"
	"github.com/zhosyaaa/RoommateTap/pkg/auth"
	"golang.org/x/oauth2"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
	cfg          oauth2.Config
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager, cfg oauth2.Config) *Handler {
	return &Handler{services: services, tokenManager: tokenManager, cfg: cfg}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		//h.initCoursesRoutes(v1)
		//h.initStudentsRoutes(v1)
		//h.initCallbackRoutes(v1)
		//h.initAdminRoutes(v1)
	}
}

//func parseIdFromPath(c *gin.Context, param string) (uint, error) {
//	idParam := c.Param(param)
//	if idParam == "" {
//		return -1, errors.New("empty id param")
//	}
//
//	id, err := (idParam)
//	if err != nil {
//		return -1, errors.New("invalid id param")
//	}
//
//	return id, nil
//}
